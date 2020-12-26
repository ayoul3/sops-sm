package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/pkg/errors"
)

func HandleDecrypt(filePath string) {
	providerClient := provider.Init()

	loader, err := GetStore(filePath)
	if err != nil {
		log.Fatal(err)
	}

	tree, err := LoadEncryptedFile(loader)
	if err != nil {
		log.Fatal(err)
	}
	if err = DecryptTree(providerClient, loader, tree); err != nil {
		log.Fatal(err)
	}
}

func LoadEncryptedFile(loader stores.StoreAPI) (*sops.Tree, error) {
	fileBytes, err := ioutil.ReadFile(loader.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}
	tree, err := loader.LoadFile(fileBytes)
	tree.FilePath = loader.GetFilePath()
	return tree, err
}

func DecryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree) (err error) {
	var content []byte
	if err = tree.Decrypt(provider); err != nil {
		return err
	}
	if content, err = loader.EmitFile(tree); err != nil {
		return err
	}
	cacheContent := tree.GetCache()
	return DumpDecryptedFiles(tree.FilePath, content, cacheContent)
}

func DumpDecryptedFiles(file string, content, cacheContent []byte) (err error) {
	cacheFile := file + ".cache"
	if err = ioutil.WriteFile(cacheFile, cacheContent, 0644); err != nil {
		return errors.Wrapf(err, "Could not write to file %s", cacheFile)
	}
	return ioutil.WriteFile(file, content, 0644)
}
