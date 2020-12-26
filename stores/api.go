package stores

import "github.com/ayoul3/sops-sm/sops"

// EncryptedFileLoader is the interface for loading of encrypted files. It provides a
// way to load encrypted SOPS files into the internal SOPS representation. Because it
// loads encrypted files, the returned data structure already contains all SOPS
// metadata.
type FileLoader interface {
	LoadFile(in []byte) (*sops.Tree, error)
}

// PlainFileEmitter is the interface for emitting plain text files. It provides a way
// to emit plain text files from the internal SOPS representation so that they can be
// shown
type FileEmitter interface {
	EmitFile(*sops.Tree) ([]byte, error)
}

// ValueEmitter is the interface for emitting a value. It provides a way to emit
// values from the internal SOPS representation so that they can be shown
type ValueEmitter interface {
	EmitValue(interface{}) ([]byte, error)
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
	ValueEmitter
	FilePathSetter
	FilePathGetter
}
