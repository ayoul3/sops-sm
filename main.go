package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ayoul3/sops-sm/lib"
	"github.com/ayoul3/sops-sm/provider/ssm"
)

func main() {
	provider := ssm.NewClient(ssm.NewAPI())
	loader, err := lib.GetStore("test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	tree, err := lib.LoadEncryptedFile(loader)
	if err != nil {
		log.Fatal(err)
	}
	out, err := lib.DecryptTree(provider, loader, tree)
	fmt.Println(out)
	fmt.Println(err)
}
