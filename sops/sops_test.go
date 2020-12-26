package sops

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var ExampleComplexTree = Tree{
	Branches: TreeBranches{
		TreeBranch{
			TreeItem{
				Key:   "hello",
				Value: `Welcome to SOPS! Edit this file as you please!`,
			},
			TreeItem{
				Key:   "example_key",
				Value: "example_value",
			},
			TreeItem{
				Key:   Comment{Value: " Example comment"},
				Value: nil,
			},
			TreeItem{
				Key: "example_array",
				Value: []interface{}{
					"example_value1",
					"example_value2",
				},
			},
			TreeItem{
				Key:   "example_number",
				Value: 1234.56789,
			},
			TreeItem{
				Key:   "example_booleans",
				Value: []interface{}{true, false},
			},
		},
	},
}

func TestSops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Sops", []Reporter{reporters.NewJUnitReporter("test_report-sops.xml")})
}

/*
var _ = Describe("LoadFile", func() {
	Context("When loading plain file succeeds", func() {
		It("should return corresponding branches", func() {
			tree, err := (&Store{}).LoadFile(PLAIN)
			Expect(err).ToNot(HaveOccurred())
			Expect(tree.Branches).To(Equal(ExampleComplexTree.Branches))
		})
	})
})
*/
