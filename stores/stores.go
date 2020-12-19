/*
Package stores acts as a layer between the internal representation of encrypted files and the encrypted files
themselves.

Subpackages implement serialization and deserialization to multiple formats.

This package defines the structure SOPS files should have and conversions to and from the internal representation. Part
of the purpose of this package is to make it easy to change the SOPS file format while remaining backwards-compatible.
*/
package stores

import (
	"time"

	"github.com/ayoul3/sops-sm/sops"
)

// SopsFile is a struct used by the stores as a helper to unmarshal the SOPS metadata
type SopsFile struct {
	// Metadata is a pointer so we can easily tell when the field is not present
	// in the SOPS file by checking for nil. This way we can show the user a
	// helpful error message indicating that the metadata wasn't found, instead
	// of showing a cryptic parsing error
	Metadata *Metadata `yaml:"sops" json:"sops" ini:"sops"`
}

type kmskey struct {
	Arn              string             `yaml:"arn" json:"arn"`
	Role             string             `yaml:"role,omitempty" json:"role,omitempty"`
	Context          map[string]*string `yaml:"context,omitempty" json:"context,omitempty"`
	CreatedAt        string             `yaml:"created_at" json:"created_at"`
	EncryptedDataKey string             `yaml:"enc" json:"enc"`
	AwsProfile       string             `yaml:"aws_profile" json:"aws_profile"`
}

type Metadata struct {
	LastModified string `yaml:"lastmodified" json:"lastmodified"`
	Version      string `yaml:"version" json:"version"`
}

// MetadataFromInternal converts an internal SOPS metadata representation to a representation appropriate for storage
func MetadataFromInternal(sopsMetadata sops.Metadata) Metadata {
	var m Metadata
	m.LastModified = sopsMetadata.LastModified.Format(time.RFC3339)
	m.Version = sopsMetadata.Version
	return m
}

// ToInternal converts a storage-appropriate Metadata struct to a SOPS internal representation
func (m *Metadata) ToInternal() (sops.Metadata, error) {
	if m == nil {
		return m.Default(), nil
	}

	lastModified, err := time.Parse(time.RFC3339, m.LastModified)
	if err != nil {
		return sops.Metadata{}, err
	}

	return sops.Metadata{
		Version:      m.Version,
		LastModified: lastModified,
	}, nil
}

func (m *Metadata) Default() sops.Metadata {
	return sops.Metadata{
		Version:      "1",
		LastModified: time.Now(),
	}
}

// ExampleComplexTree is an example sops.Tree object exhibiting complex relationships
var ExampleComplexTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   "hello",
				Value: `Welcome to SOPS! Edit this file as you please!`,
			},
			sops.TreeItem{
				Key:   "example_key",
				Value: "example_value",
			},
			sops.TreeItem{
				Key:   sops.Comment{Value: " Example comment"},
				Value: nil,
			},
			sops.TreeItem{
				Key: "example_array",
				Value: []interface{}{
					"example_value1",
					"example_value2",
				},
			},
			sops.TreeItem{
				Key:   "example_number",
				Value: 1234.56789,
			},
			sops.TreeItem{
				Key:   "example_booleans",
				Value: []interface{}{true, false},
			},
		},
	},
}

// ExampleSimpleTree is an example sops.Tree object exhibiting only simple relationships
// with only one nested branch and only simple string values
var ExampleSimpleTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "Welcome!",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   sops.Comment{Value: " This is an example file."},
						Value: nil,
					},
					sops.TreeItem{
						Key:   "hello",
						Value: "Welcome to SOPS! Edit this file as you please!",
					},
					sops.TreeItem{
						Key:   "example_key",
						Value: "example_value",
					},
				},
			},
		},
	},
}

// ExampleFlatTree is an example sops.Tree object exhibiting only simple relationships
// with no nested branches and only simple string values
var ExampleFlatTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   sops.Comment{Value: " This is an example file."},
				Value: nil,
			},
			sops.TreeItem{
				Key:   "hello",
				Value: "Welcome to SOPS! Edit this file as you please!",
			},
			sops.TreeItem{
				Key:   "example_key",
				Value: "example_value",
			},
			sops.TreeItem{
				Key:   "example_multiline",
				Value: "foo\nbar\nbaz",
			},
		},
	},
}
