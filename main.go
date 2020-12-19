package main

import (
	"fmt"

	"github.com/ayoul3/sops-sm/lib"
)

func main() {
	loader, err := lib.GetStore("text.yaml")
	fmt.Println(err)
	tree, err := lib.LoadEncryptedFile(loader, "test.yaml")
	fmt.Println(err)
	fmt.Println(tree)
}
