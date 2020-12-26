package yaml

import (
	"testing"

	"github.com/ayoul3/sops-sm/sops"
	"github.com/mozilla-services/yaml"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var PLAIN = []byte(`---
hello: Welcome to SOPS! Edit this file as you please!
example_key: example_value
# Example comment
example_array:
- example_value1
- example_value2
example_number: 1234.56789
example_booleans:
- true
- false`)

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

func TestYAML(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Yaml", []Reporter{reporters.NewJUnitReporter("test_report-yaml.xml")})
}

var _ = Describe("LoadFile", func() {
	Context("When loading plain file succeeds", func() {
		It("should return corresponding branches", func() {
			tree, err := (&Store{}).LoadFile(PLAIN)
			Expect(err).ToNot(HaveOccurred())
			Expect(tree.Branches).To(Equal(ExampleComplexTree.Branches))
		})
	})
	Context("When loading plain file fails", func() {
		It("should return an error", func() {
			_, err := (&Store{}).LoadFile([]byte(`---\nkey1: va:lue\n:`))
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = Describe("EmitFile", func() {
	Context("When loading a tree succeeds", func() {
		It("should return a yaml file", func() {
			var generic interface{}
			content, err := (&Store{}).EmitFile(&ExampleComplexTree)
			Expect(err).ToNot(HaveOccurred())
			err = yaml.Unmarshal(content, &generic)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(ContainSubstring("example_key: example_value"))
			Expect(content).To(ContainSubstring("# Example comment"))
		})
	})
})
