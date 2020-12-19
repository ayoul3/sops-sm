package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ayoul3/sops-sm/sops"
)

func LoadEncryptedFile(loader sops.Store, inputPath string) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}
	path, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, err
	}
	tree, err := loader.LoadEncryptedFile(fileBytes)
	tree.FilePath = path
	return &tree, err
}
