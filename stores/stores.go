package stores

import (
	"github.com/ayoul3/sops-sm/sops"
)

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
