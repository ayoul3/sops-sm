package cmd

import (
	"io/ioutil"
	"log"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/pkg/errors"
)

func HandleEncrypt(filePath string) {
	providerClient := provider.Init()

	loader, err := GetStore(filePath)
	if err != nil {
		log.Fatal(err)
	}

	tree, err := LoadPlainFile(loader)
	if err != nil {
		log.Fatal(err)
	}

	if err = EncryptTree(providerClient, loader, tree); err != nil {
		log.Fatal(err)
	}
}

func LoadPlainFile(loader stores.StoreAPI) (outTree *sops.Tree, err error) {
	var tree sops.Tree
	fileBytes, err := ioutil.ReadFile(loader.GetFilePath())
	if err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error reading file ")
	}
	if tree, err = loader.LoadFile(fileBytes); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error loading file ")
	}
	if err = tree.LoadCache(loader.GetCachePath()); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error loading cache file ")
	}
	tree.FilePath = loader.GetFilePath()
	return &tree, nil
}

func EncryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree) (err error) {
	var content []byte
	if err = tree.Encrypt(provider); err != nil {
		return err
	}
	if content, err = loader.EmitPlainFile(tree.Branches); err != nil {
		return err
	}
	return DumpPlainFile(tree.FilePath, content)
}

func DumpPlainFile(file string, content []byte) (err error) {
	return ioutil.WriteFile(file, content, 0644)
}
