package stores

import "github.com/ayoul3/sops-sm/sops"

// Store is used to interact with files, both encrypted and unencrypted.
type StoreAPI interface {
	LoadFile(in []byte) (*sops.Tree, error)
	EmitFile(*sops.Tree) ([]byte, error)
	SetFilePath(string)
	GetFilePath() string
	GetCachePath() string
}
