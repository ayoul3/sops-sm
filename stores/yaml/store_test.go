package yaml

import (
	"testing"

	"github.com/ayoul3/sops-sm/sops"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var PLAIN = []byte(`---
# comment 0
key1: value
key1_a: value
# ^ comment 1
---
key2: value2`)

var BRANCHES = sops.TreeBranches{
	sops.TreeBranch{
		sops.TreeItem{
			Key:   "key1",
			Value: "value",
		},
		sops.TreeItem{
			Key:   "key1_a",
			Value: "value",
		},
		sops.TreeItem{
			Key:   sops.Comment{" ^ comment 1"},
			Value: nil,
		},
	},
	sops.TreeBranch{
		sops.TreeItem{
			Key:   "key2",
			Value: "value2",
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
			Expect(tree.Branches).To(Equal(BRANCHES))
		})
	})
	Context("When loading plain file fails", func() {
		It("should return an error", func() {
			_, err := (&Store{}).LoadFile([]byte(`---\nkey1: va:lue\n:`))
			Expect(err).To(HaveOccurred())
		})
	})
})
