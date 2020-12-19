package sops

import "time"

// Comment represents a comment in the sops tree for the file formats that actually support them.
type Comment struct {
	Value string
}

// TreeItem is an item inside sops's tree
type TreeItem struct {
	Key   interface{}
	Value interface{}
}

// TreeBranch is a branch inside sops's tree. It is a slice of TreeItems and is therefore ordered
type TreeBranch []TreeItem

// TreeBranches is a collection of TreeBranch
// Trees usually have more than one branch
type TreeBranches []TreeBranch

// Tree is the data structure used by sops to represent documents internally
type Tree struct {
	Metadata Metadata
	Branches TreeBranches
	// FilePath is the path of the file this struct represents
	FilePath string
}

// Metadata holds information about a file encrypted by sops
type Metadata struct {
	LastModified time.Time
	Version      string
	// DataKey caches the decrypted data key so it doesn't have to be decrypted with a master key every time it's needed
	DataKey []byte
}

// EncryptedFileLoader is the interface for loading of encrypted files. It provides a
// way to load encrypted SOPS files into the internal SOPS representation. Because it
// loads encrypted files, the returned data structure already contains all SOPS
// metadata.
type EncryptedFileLoader interface {
	LoadEncryptedFile(in []byte) (Tree, error)
}

// PlainFileLoader is the interface for loading of plain text files. It provides a
// way to load unencrypted files into SOPS. Because the files it loads are
// unencrypted, the returned data structure does not contain any metadata.
type PlainFileLoader interface {
	LoadPlainFile(in []byte) (TreeBranches, error)
}

// EncryptedFileEmitter is the interface for emitting encrypting files. It provides a
// way to emit encrypted files from the internal SOPS representation.
type EncryptedFileEmitter interface {
	EmitEncryptedFile(Tree) ([]byte, error)
}

// PlainFileEmitter is the interface for emitting plain text files. It provides a way
// to emit plain text files from the internal SOPS representation so that they can be
// shown
type PlainFileEmitter interface {
	EmitPlainFile(TreeBranches) ([]byte, error)
}

// ValueEmitter is the interface for emitting a value. It provides a way to emit
// values from the internal SOPS representation so that they can be shown
type ValueEmitter interface {
	EmitValue(interface{}) ([]byte, error)
}

// Store is used to interact with files, both encrypted and unencrypted.
type Store interface {
	EncryptedFileLoader
	PlainFileLoader
	EncryptedFileEmitter
	PlainFileEmitter
	ValueEmitter
}
