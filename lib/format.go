package lib

import (
	"fmt"

	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores/yaml"
)

func newYamlStore() sops.Store {
	return &yaml.Store{}
}

var stores = map[string]sops.Store{
	//Json: newJsonStore,
	"yaml": newYamlStore(),
}

func getFileFormat(inputFile string) string {
	return "yaml"

}
func GetStore(inputFile string) (sops.Store, error) {
	format := getFileFormat(inputFile)
	if val, ok := stores[format]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("File format not supported: %s", format)
}
