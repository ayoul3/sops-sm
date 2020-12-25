package sops

import (
	"fmt"
	"reflect"
	"strings"
	"time"

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

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Metadata Metadata
	Branches TreeBranches
	// FilePath is the path of the file this struct represents
	FilePath string
}

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified time.Time
	Version      string
	// DataKey caches the decrypted data key so it doesn't have to be decrypted with a master key every time it's needed
	DataKey map[string]string
}

func valueFromPathAndLeaf(path []interface{}, leaf interface{}) interface{} {
	switch component := path[0].(type) {
	case int:
		if len(path) == 1 {
			return []interface{}{
				leaf,
			}
		}
		return []interface{}{
			valueFromPathAndLeaf(path[1:], leaf),
		}
	default:
		if len(path) == 1 {
			return TreeBranch{
				TreeItem{
					Key:   component,
					Value: leaf,
				},
			}
		}
		return TreeBranch{
			TreeItem{
				Key:   component,
				Value: valueFromPathAndLeaf(path[1:], leaf),
			},
		}
	}
}

func set(branch interface{}, path []interface{}, value interface{}) interface{} {
	switch branch := branch.(type) {
	case TreeBranch:
		for i, item := range branch {
			if item.Key == path[0] {
				if len(path) == 1 {
					branch[i].Value = value
				} else {
					branch[i].Value = set(item.Value, path[1:], value)
				}
				return branch
			}
		}
		// Not found, need to add the next path entry to the branch
		if len(path) == 1 {
			return append(branch, TreeItem{Key: path[0], Value: value})
		}
		return valueFromPathAndLeaf(path, value)
	case []interface{}:
		position := path[0].(int)
		if len(path) == 1 {
			if position >= len(branch) {
				return append(branch, value)
			}
			branch[position] = value
		} else {
			if position >= len(branch) {
				branch = append(branch, valueFromPathAndLeaf(path[1:], value))
			}
			branch[position] = set(branch[position], path[1:], value)
		}
		return branch
	default:
		return valueFromPathAndLeaf(path, value)
	}
}

// Set sets a value on a given tree for the specified path
func (branch TreeBranch) Set(path []interface{}, value interface{}) TreeBranch {
	return set(branch, path, value).(TreeBranch)
}

// Truncate truncates the tree to the path specified
func (branch TreeBranch) Truncate(path []interface{}) (interface{}, error) {
	log.WithField("path", path).Info("Truncating tree")
	var current interface{} = branch
	for _, component := range path {
		switch component := component.(type) {
		case string:
			found := false
			for _, item := range current.(TreeBranch) {
				if item.Key == component {
					current = item.Value
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("component ['%s'] not found", component)
			}
		case int:
			if reflect.ValueOf(current).Kind() != reflect.Slice {
				return nil, fmt.Errorf("component [%d] is integer, but tree part is not a slice", component)
			}
			if reflect.ValueOf(current).Len() <= component {
				return nil, fmt.Errorf("component [%d] accesses out of bounds", component)
			}
			current = reflect.ValueOf(current).Index(component).Interface()
		}
	}
	return current, nil
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

// Decrypt walks over the tree and decrypts all values with the provided cipher,
// except those whose key ends with the UnencryptedSuffix specified on the Metadata struct,
// those not ending with EncryptedSuffix, if EncryptedSuffix is provided (by default it is not),
// those not matching EncryptedRegex, if EncryptedRegex is provided (by default it is not),
// or those matching UnencryptedRegex, if UnencryptedRegex is provided (by default it is not).
// If decryption is successful, it returns the MAC for the decrypted tree.
func (tree Tree) Decrypt(provider provider.API) error {
	log.Info("Decrypting tree")
	tree.Metadata.DataKey = make(map[string]string, 0)

	walk := func(branch TreeBranch) error {
		_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (interface{}, error) {
			var v, cached string
			var ok, found bool

			pathString := strings.Join(path, ":") + ":"
			log.Infof("Walking path %s ", pathString)

			if v, ok = in.(string); !ok {
				return nil, fmt.Errorf("Decrypt - expecting string value for secret key ID")
			}
			if cached, found = tree.Metadata.DataKey[v]; found {
				log.Infof("Found secret in cache %s", v)
				return cached, nil
			}
			if provider.IsSecret(v) {
				log.Infof("Fetching secret %s ", v)
				secretValue, err := provider.GetSecret(v)
				tree.Metadata.DataKey[v] = secretValue
				return secretValue, err
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
