package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ayoul3/sops-sm/stores"
	"github.com/ayoul3/sops-sm/stores/json"
	"github.com/ayoul3/sops-sm/stores/yaml"
	"github.com/spf13/afero"
)

var formats = map[string]stores.StoreAPI{
	"yaml": yaml.NewStore(),
	"json": json.NewStore(),
}

type Handler struct {
	Fs         afero.Fs
	NumThreads int
	Overwrite  bool
}

func NewHandler(numThreads int, overwrite bool) *Handler {
	return &Handler{
		Fs:         afero.NewOsFs(),
		NumThreads: numThreads,
		Overwrite:  overwrite,
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

func (h *Handler) GetOutputFileName(filePath, suffix string) string {
	if h.Overwrite {
		return filePath
	}
	fileParts := strings.Split(filePath, ".")
	ext := filepath.Ext(filePath)
	if len(fileParts) < 1 {
		return filePath + suffix
	}
	return strings.Join(fileParts[0:len(fileParts)-1], "") + suffix + ext
}
