package lib

import (
	"fmt"
	"path/filepath"

	"github.com/ayoul3/sops-sm/stores"
	"github.com/ayoul3/sops-sm/stores/json"
	"github.com/ayoul3/sops-sm/stores/yaml"
)

var formats = map[string]stores.StoreAPI{
	"yaml": yaml.NewStore(),
	"json": json.NewStore(),
}

func getFileFormat(inputFile string) string {
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

func GetStore(inputFile string) (stores.StoreAPI, error) {
	format := getFileFormat(inputFile)
	if val, ok := formats[format]; ok {
		val.SetFilePath(inputFile)
		return val, nil
	}
	return nil, fmt.Errorf("File format not supported: %s", format)
}
