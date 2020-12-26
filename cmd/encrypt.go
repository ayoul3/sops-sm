package cmd

import (
	"log"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func (h *Handler) HandleEncrypt(filePath string) {
	var content []byte
	providerClient := provider.Init()

	loader, err := h.GetStore(filePath)
	if err != nil {
		log.Fatal(err)
	}

	tree, err := LoadPlainFile(h, loader)
	if err != nil {
		log.Fatal(err)
	}

	if content, err = EncryptTree(providerClient, loader, tree); err != nil {
		log.Fatal(err)
	}

	if err = DumpPlainFile(h, tree.FilePath, content); err != nil {
		log.Fatal(err)
	}
}

func LoadPlainFile(h *Handler, loader stores.StoreAPI) (tree *sops.Tree, err error) {
	var fileBytes []byte
	var cacheReader afero.File

	if fileBytes, err = afero.ReadFile(h.Fs, loader.GetFilePath()); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error reading file ")
	}
	if tree, err = loader.LoadFile(fileBytes); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error loading file ")
	}

	if cacheReader, err = h.Fs.Open(loader.GetCachePath()); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error reading file ")
	}
	defer cacheReader.Close()

	if err = tree.LoadCache(cacheReader); err != nil {
		return nil, errors.Wrap(err, "LoadPlainFile: Error loading cache file ")
	}
	tree.FilePath = loader.GetFilePath()
	return tree, nil
}

func EncryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree) (content []byte, err error) {
	if err = tree.Encrypt(provider); err != nil {
		return
	}
	if content, err = loader.EmitFile(tree); err != nil {
		return
	}
	return
}

func DumpPlainFile(h *Handler, file string, content []byte) (err error) {
	return afero.WriteFile(h.Fs, file, content, 0644)
}
