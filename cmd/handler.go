package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ayoul3/sops-sm/stores"
	"github.com/spf13/afero"
)

type Handler struct {
	Fs         afero.Fs
	numThreads int
}

func NewHandler(numThreads int) *Handler {
	return &Handler{
		Fs:         afero.NewOsFs(),
		numThreads: numThreads,
	}
}

func GetFileFormat(inputFile string) string {
	extension := filepath.Ext(inputFile)
	if len(extension) < 2 {
		return ""
	}
	extension = extension[1:]
	switch extension {
	case "yaml", "yml":
		return "yaml"
	default:
		return extension
	}
}

func (h *Handler) GetStore(inputFile string) (stores.StoreAPI, error) {
	format := GetFileFormat(inputFile)
	if val, ok := formats[format]; ok {
		val.SetFilePath(inputFile)
		return val, nil
	}
	return nil, fmt.Errorf("File format not supported: %s", format)
}
