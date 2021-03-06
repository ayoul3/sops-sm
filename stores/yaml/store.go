package yaml //import "go.mozilla.org/sops/v3/stores/yaml"

import (
	"fmt"

	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/mozilla-services/yaml"
)

// Store handles storage of YAML data
type Store struct {
	path string
}

func NewStore() stores.StoreAPI {
	return &Store{}
}

func (store *Store) GetFilePath() string {
	return store.path
}

func (store *Store) GetCachePath() string {
	return store.path + ".cache"
}

func (store *Store) SetFilePath(p string) {
	store.path = p
}

func (store Store) mapSliceToTreeBranch(in yaml.MapSlice) sops.TreeBranch {
	branch := make(sops.TreeBranch, 0)
	for _, item := range in {
		if comment, ok := item.Key.(yaml.Comment); ok {
			// Convert the yaml comment to a generic sops comment
			branch = append(branch, sops.TreeItem{
				Key: sops.Comment{
					Value: comment.Value,
				},
				Value: nil,
			})
		} else {
			branch = append(branch, sops.TreeItem{
				Key:   item.Key,
				Value: store.yamlValueToTreeValue(item.Value),
			})
		}
	}
	return branch
}

func (store Store) yamlValueToTreeValue(in interface{}) interface{} {
	switch in := in.(type) {
	case map[interface{}]interface{}:
		return store.yamlMapToTreeBranch(in)
	case yaml.MapSlice:
		return store.mapSliceToTreeBranch(in)
	case []interface{}:
		return store.yamlSliceToTreeValue(in)
	case yaml.Comment:
		return sops.Comment{Value: in.Value}
	default:
		return in
	}
}

func (store *Store) yamlSliceToTreeValue(in []interface{}) []interface{} {
	for i, v := range in {
		in[i] = store.yamlValueToTreeValue(v)
	}
	return in
}

func (store *Store) yamlMapToTreeBranch(in map[interface{}]interface{}) sops.TreeBranch {
	branch := make(sops.TreeBranch, 0)
	for k, v := range in {
		branch = append(branch, sops.TreeItem{
			Key:   k.(string),
			Value: store.yamlValueToTreeValue(v),
		})
	}
	return branch
}

func (store Store) treeValueToYamlValue(in interface{}) interface{} {
	switch in := in.(type) {
	case sops.TreeBranch:
		return store.treeBranchToYamlMap(in)
	case sops.Comment:
		return yaml.Comment{in.Value}
	case []interface{}:
		var out []interface{}
		for _, v := range in {
			out = append(out, store.treeValueToYamlValue(v))
		}
		return out
	default:
		return in
	}
}

func (store Store) treeBranchToYamlMap(in sops.TreeBranch) yaml.MapSlice {
	branch := make(yaml.MapSlice, 0)
	for _, item := range in {
		if comment, ok := item.Key.(sops.Comment); ok {
			branch = append(branch, yaml.MapItem{
				Key:   store.treeValueToYamlValue(comment),
				Value: nil,
			})
		} else {
			branch = append(branch, yaml.MapItem{
				Key:   item.Key,
				Value: store.treeValueToYamlValue(item.Value),
			})
		}
	}
	return branch
}

func (store *Store) LoadFile(in []byte) (*sops.Tree, error) {
	var data []yaml.MapSlice
	if err := (yaml.CommentUnmarshaler{}).UnmarshalDocuments(in, &data); err != nil {
		return &sops.Tree{}, fmt.Errorf("Error unmarshaling input YAML: %s", err)
	}
	var branches sops.TreeBranches
	for _, doc := range data {
		branches = append(branches, store.mapSliceToTreeBranch(doc))
	}
	return &sops.Tree{
		Branches: branches,
		Cache:    make(map[string]sops.CachedSecret, 0),
	}, nil
}

// EmitFile returns the encrypted bytes of the yaml file corresponding to a
// sops.Tree runtime object
func (store *Store) EmitFile(in *sops.Tree) ([]byte, error) {
	out := []byte{}
	for i, branch := range in.Branches {
		if i > 0 {
			out = append(out, "---\n"...)
		}
		yamlMap := store.treeBranchToYamlMap(branch)
		tout, err := (&yaml.YAMLMarshaler{Indent: 4}).Marshal(yamlMap)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling to yaml: %s", err)
		}
		out = append(out, tout...)
	}
	return out, nil
}
