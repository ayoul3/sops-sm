package stores

import "github.com/ayoul3/sops-sm/sops"

// FileLoader is the interface for loading files
type FileLoader interface {
	LoadFile(in []byte) (*sops.Tree, error)
}

// FileEmitter is the interface for emitting files.
type FileEmitter interface {
	EmitFile(*sops.Tree) ([]byte, error)
}

// FilePathSetter is the interface for setting the filepath of the loaded file
type FilePathSetter interface {
	SetFilePath(string)
}

// FilePathGetter is the interface for getting the filepath of the loaded file
type FilePathGetter interface {
	GetFilePath() string
	GetCachePath() string
}

// Store is used to interact with files, both encrypted and unencrypted.
type StoreAPI interface {
	FileLoader
	FileEmitter
	FilePathSetter
	FilePathGetter
}
