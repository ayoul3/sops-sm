package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
)

func LoadEncryptedFile(loader stores.StoreAPI) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(loader.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}
	path, err := filepath.Abs(loader.GetFilePath())
	if err != nil {
		return nil, err
	}
	tree, err := loader.LoadEncryptedFile(fileBytes)
	tree.FilePath = path
	return &tree, err
}

func DecryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree) (decryptedFile []byte, err error) {
	fmt.Println(tree)
	err = tree.Decrypt(provider)
	fmt.Println(tree)
	return nil, err
}
