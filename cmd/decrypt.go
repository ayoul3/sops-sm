package cmd

import (
	"fmt"
	"log"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/sops"
	"github.com/ayoul3/sops-sm/stores"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func (h *Handler) HandleDecrypt(filePath string) {
	var content []byte

	providerClient := provider.Init()
	loader, err := h.GetStore(filePath)
	if err != nil {
		log.Fatal(err)
	}
	tree, err := LoadEncryptedFile(h, loader)
	if err != nil {
		log.Fatal(err)
	}
	if content, err = DecryptTree(providerClient, loader, tree, h.numThreads); err != nil {
		log.Fatal(err)
	}
	if err = DumpDecryptedTree(h, tree.FilePath, loader.GetCachePath(), content, tree.GetCache()); err != nil {
		log.Fatal(err)
	}
}

func LoadEncryptedFile(h *Handler, loader stores.StoreAPI) (*sops.Tree, error) {
	fileBytes, err := afero.ReadFile(h.Fs, loader.GetFilePath())
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err)
	}
	tree, err := loader.LoadFile(fileBytes)
	tree.FilePath = loader.GetFilePath()
	return tree, err
}

func DecryptTree(provider provider.API, loader stores.StoreAPI, tree *sops.Tree, numThreads int) (content []byte, err error) {
	if err = tree.Decrypt(provider, numThreads); err != nil {
		return
	}
	if content, err = loader.EmitFile(tree); err != nil {
		return
	}
	return
}

func DumpDecryptedTree(h *Handler, file, cacheFile string, content, cacheContent []byte) (err error) {
	if err = afero.WriteFile(h.Fs, cacheFile, cacheContent, 0644); err != nil {
		return errors.Wrapf(err, "Could not write to file %s", cacheFile)
	}
	return afero.WriteFile(h.Fs, file, content, 0644)
}
