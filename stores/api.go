package stores

import "github.com/ayoul3/sops-sm/sops"

// EncryptedFileLoader is the interface for loading of encrypted files. It provides a
// way to load encrypted SOPS files into the internal SOPS representation. Because it
// loads encrypted files, the returned data structure already contains all SOPS
// metadata.
type EncryptedFileLoader interface {
	LoadEncryptedFile(in []byte) (sops.Tree, error)
}

// PlainFileLoader is the interface for loading of plain text files. It provides a
// way to load unencrypted files into SOPS. Because the files it loads are
// unencrypted, the returned data structure does not contain any metadata.
type PlainFileLoader interface {
	LoadPlainFile(in []byte) (sops.TreeBranches, error)
}

// EncryptedFileEmitter is the interface for emitting encrypting files. It provides a
// way to emit encrypted files from the internal SOPS representation.
type EncryptedFileEmitter interface {
	EmitEncryptedFile(sops.Tree) ([]byte, error)
}

// PlainFileEmitter is the interface for emitting plain text files. It provides a way
// to emit plain text files from the internal SOPS representation so that they can be
// shown
type PlainFileEmitter interface {
	EmitPlainFile(sops.TreeBranches) ([]byte, error)
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
}

// Store is used to interact with files, both encrypted and unencrypted.
type StoreAPI interface {
	EncryptedFileLoader
	PlainFileLoader
	EncryptedFileEmitter
	PlainFileEmitter
	ValueEmitter
	FilePathSetter
	FilePathGetter
}
