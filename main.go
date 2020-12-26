package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ayoul3/sops-sm/lib"
	"github.com/ayoul3/sops-sm/provider"
)

var decrypt, encrypt bool
var filePath string

func init() {
	flag.BoolVar(&decrypt, "d", true, "Decode the input file")
	flag.BoolVar(&encrypt, "e", false, "Encode the input file - need .cache file genrated by decoding process")
	flag.Parse()
	ValidateParams()
}

func main() {
	providerClient := provider.Init()
	loader, err := lib.GetStore(filePath)
	if err != nil {
		log.Fatal(err)
	}

	tree, err := lib.LoadEncryptedFile(loader)
	if err != nil {
		log.Fatal(err)
	}
	if err = lib.DecryptTree(providerClient, loader, tree); err != nil {
		log.Fatal(err)
	}
}

func ValidateParams() {
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filePath = flag.Arg(0)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		flag.Usage()
		log.Fatalf("input file %s does not exist", filePath)
	}
}
