package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/pkg/errors"
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

func DecryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree) (err error) {
	var content []byte
	if err = tree.Decrypt(provider); err != nil {
		return err
	}
	if content, err = loader.EmitPlainFile(tree.Branches); err != nil {
		return err
	}
	cacheContent := tree.GetCache()
	return DumpFiles(tree.FilePath, content, cacheContent)
}

func DumpFiles(file string, content, cacheContent []byte) (err error) {
	cacheFile := file + ".cache"
	if err = ioutil.WriteFile(cacheFile, cacheContent, 0644); err != nil {
		return errors.Wrapf(err, "Could not write to file %s", cacheFile)
	}
	return ioutil.WriteFile(file, content, 0644)
}
