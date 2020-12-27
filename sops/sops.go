package sops

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ayoul3/sops-sm/provider"
	log "github.com/sirupsen/logrus"
)

// Comment represents a comment in the sops tree for the file formats that actually support them.
type Comment struct {
	Value string
}

// TreeItem is an item inside sops's tree
type TreeItem struct {
	Key   interface{}
	Value interface{}
}

// TreeBranch is a branch inside sops's tree. It is a slice of TreeItems and is therefore ordered
type TreeBranch []TreeItem

// TreeBranches is a collection of TreeBranch
// Trees usually have more than one branch
type TreeBranches []TreeBranch

type PathSecret struct {
	FullKey  string
	FullPath string
}

type CachedSecret struct {
	Value string
	Path  []PathSecret
}

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Cache    map[string]CachedSecret
	Branches TreeBranches
	// FilePath is the path of the file this struct represents
	FilePath string
}

func (branch TreeBranch) walkValue(in interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (interface{}, error) {
	switch in := in.(type) {
	case string:
		return onLeaves(in, path)
	case []byte:
		return onLeaves(string(in), path)
	case int:
		return onLeaves(in, path)
	case bool:
		return onLeaves(in, path)
	case float64:
		return onLeaves(in, path)
	case Comment:
		return onLeaves(in, path)
	case TreeBranch:
		return branch.walkBranch(in, path, onLeaves)
	case []interface{}:
		return branch.walkSlice(in, path, onLeaves)
	case nil:
		// the value returned remains the same since it doesn't make
		// sense to encrypt or decrypt a nil value
		return nil, nil
	default:
		return nil, fmt.Errorf("Cannot walk value, unknown type: %T", in)
	}
}

func (branch TreeBranch) walkSlice(in []interface{}, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) ([]interface{}, error) {
	for i, v := range in {
		newV, err := branch.walkValue(v, path, onLeaves)
		if err != nil {
			return nil, err
		}
		in[i] = newV
	}
	return in, nil
}

func (branch TreeBranch) walkBranch(in TreeBranch, path []string, onLeaves func(in interface{}, path []string) (interface{}, error)) (TreeBranch, error) {
	for i, item := range in {
		if _, ok := item.Key.(Comment); ok {
			enc, err := branch.walkValue(item.Key, path, onLeaves)
			if err != nil {
				return nil, err
			}
			if encComment, ok := enc.(Comment); ok {
				in[i].Key = encComment
				continue
			} else if comment, ok := enc.(string); ok {
				in[i].Key = Comment{Value: comment}
				continue
			} else {
				return nil, fmt.Errorf("walkValue of Comment should be either Comment or string, was %T", enc)
			}
		}
		key, ok := item.Key.(string)
		if !ok {
			return nil, fmt.Errorf("Tree contains a non-string key (type %T): %s. Only string keys are"+
				"supported", item.Key, item.Key)
		}
		newV, err := branch.walkValue(item.Value, append(path, key), onLeaves)
		if err != nil {
			return nil, err
		}
		in[i].Value = newV
	}
	return in, nil
}

// Decrypt walks over the tree and fetches IDs from SecretsManager or ParameterStore
func (tree *Tree) DecryptAsync(provider provider.API) (err error) {
	log.Info("First walk down the tree to fetch secrets")
	for _, branch := range tree.Branches {
		if err = WalkerAsyncFetchSecret(branch, provider); err != nil {
			return fmt.Errorf("Error walking tree: %s", err)
		}
	}
	CacheAsyncSecret(tree)
	log.Info("Second walk down the tree to fetch secrets from cache")
	for _, branch := range tree.Branches {
		if err = WalkerSyncFetchSecret(tree, branch, provider); err != nil {
			return fmt.Errorf("Error walking tree: %s", err)
		}
	}
	return nil
}

// Decrypt walks over the tree and fetches IDs from SecretsManager or ParameterStore
func (tree *Tree) DecryptSync(provider provider.API) (err error) {
	log.Info("Decrypting tree")

	for _, branch := range tree.Branches {
		if err = WalkerSyncFetchSecret(tree, branch, provider); err != nil {
			return fmt.Errorf("Error walking tree: %s", err)
		}
	}
	return nil
}

func (tree *Tree) Decrypt(provider provider.API, numThreads int) (err error) {
	if numThreads > 1 {
		PrepareAsync(tree, provider, numThreads)
		return tree.DecryptAsync(provider)
	}
	return tree.DecryptSync(provider)
}

// Decrypt walks over the tree and fetches IDs from SecretsManager or ParameterStore
func (tree Tree) Encrypt(provider provider.API) error {
	log.Info("Encrypting tree")

	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (v interface{}, err error) {
			var cached CachedSecret
			var ok, found bool

			pathString := strings.Join(path, ":")
			log.Infof("Walking path %s ", pathString)

			if v, ok = in.(string); !ok {
				return in, nil
			}
			if cached, found = tree.Cache[pathString]; found {
				log.Infof("Found secret in cache %s", v)
				return cached.Value, nil
			}
			return v, nil

		})
		return err
	}
	for _, branch := range tree.Branches {
		err := walk(branch)
		if err != nil {
			return fmt.Errorf("Error walking tree: %s", err)
		}
	}
	return nil
}

func (tree Tree) GetCache() []byte {
	out := bytes.NewBuffer([]byte(""))
	for _, secret := range tree.Cache {
		for _, path := range secret.Path {
			c := fmt.Sprintf("%s,%s\n", path.FullPath, path.FullKey)
			out.WriteString(c)
		}
	}
	return out.Bytes()
}

func (tree *Tree) CacheSecretValue(fullKey, value, path string) {
	var re = regexp.MustCompile(`@.+`)
	baseKey := re.ReplaceAllString(fullKey, ``)

	pathSecret := []PathSecret{{FullKey: fullKey, FullPath: path}}
	if stored, ok := tree.Cache[baseKey]; ok {
		pathSecret = append(stored.Path, pathSecret...)
	}
	tree.Cache[baseKey] = CachedSecret{Value: value, Path: pathSecret}
	return
}

func (tree *Tree) IsCached(key string) (string, bool) {
	var re = regexp.MustCompile(`@.+`)
	key = re.ReplaceAllString(key, ``)
	secret, found := tree.Cache[key]
	return secret.Value, found
}

func (tree *Tree) LoadCache(fileReader io.Reader) {
	tree.Cache = make(map[string]CachedSecret, 0)
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		line := strings.Split(strings.Trim(scanner.Text(), "\n"), ",")
		if len(line) < 2 {
			log.Warnf("LoadCache: ignoring line %s", line)
			continue
		}
		tree.Cache[line[0]] = CachedSecret{Value: line[1]}
	}
}
