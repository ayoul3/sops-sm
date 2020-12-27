package sops

import (
	"testing"

	"github.com/ayoul3/sops-sm/provider/ssm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func getTree() Tree {
	return Tree{
		Cache: make(map[string]CachedSecret, 0),
		Branches: TreeBranches{
			TreeBranch{
				TreeItem{
					Key: "hello",
					Value: TreeBranch{
						TreeItem{
							Key:   "nested",
							Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam",
						},
					},
				},
				TreeItem{
					Key:   "secret",
					Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam",
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
					Key:   "example_int",
					Value: 1234,
				},
				TreeItem{
					Key:   "test",
					Value: nil,
				},
				TreeItem{
					Key:   "example_booleans",
					Value: []interface{}{true, false},
				},
			},
		},
	}
}

func TestSops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Sops", []Reporter{reporters.NewJUnitReporter("test_report-sops.xml")})
}

var _ = Describe("DecryptSync", func() {
	Context("When decrypting a tree succeeds", func() {
		It("should return branches with secret value", func() {
			client := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			tree := getTree()
			err := tree.DecryptSync(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(tree.Branches[0][0].Value).To(Equal(TreeBranch{TreeItem{Key: "nested", Value: "test"}}))
			Expect(tree.Branches[0][1].Value).To(Equal("test"))
			Expect(len(tree.Cache)).To(Equal(1))
			Expect(tree.Cache).To(HaveKey("arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"))
			Expect(tree.Cache["arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"].Value).To(Equal("test"))
			Expect(len(tree.Cache["arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"].Path)).To(Equal(2))
		})
	})
	Context("When decrypting a tree with index", func() {
		It("should return branches with secret value", func() {
			client := ssm.NewClient(&ssm.MockClient{SecretValue: `{"index1": "value1", "index2": "value2"}`})
			tree := getTree()
			tree.Branches[0][1].Value = "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam@index1"
			tree.Branches[0][3].Value = []interface{}{
				"arn:aws:ssm:eu-west-1:123456789123:parameter/someparam@index2",
				"example_value2",
			}

			err := tree.DecryptSync(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(tree.Branches[0][1].Value).To(Equal("value1"))
			Expect(tree.Branches[0][3].Value).To(Equal([]interface{}{
				"value2",
				"example_value2",
			}))
			Expect(len(tree.Cache)).To(Equal(1))
			Expect(tree.Cache).To(HaveKey("arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"))
			Expect(tree.Cache["arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"].Value).To(Equal(`{"index1": "value1", "index2": "value2"}`))
			Expect(len(tree.Cache["arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"].Path)).To(Equal(3))
		})
	})
	Context("When decrypting a tree fails", func() {
		It("should return an error", func() {
			client := ssm.NewClient(&ssm.MockClient{GetParameterShouldFail: true})
			tree := getTree()
			err := tree.DecryptSync(client)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Encrypt", func() {
	Context("When encrypting a tree succeeds", func() {
		It("should return branches with secret value", func() {
			client := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			tree := getTree()
			tree.Branches[0][0].Value = "test"
			tree.Branches[0][0].Value = TreeBranch{TreeItem{Key: "nested", Value: "test"}}
			tree.Cache = map[string]CachedSecret{
				"hello:nested": {Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"},
				"secret":       {Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"},
			}
			err := tree.Encrypt(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(tree.Branches[0][0].Value).To(Equal(TreeBranch{TreeItem{Key: "nested", Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"}}))
			Expect(tree.Branches[0][1].Value).To(Equal("arn:aws:ssm:eu-west-1:123456789123:parameter/someparam"))
		})
	})
})

var _ = Describe("LoadCache", func() {
	Context("When loading the cache succeeds", func() {
		It("should return a full cache", func() {
			fs := afero.NewMemMapFs()
			afero.WriteFile(fs, "test.yaml.cache", []byte("values:key1,arn:aws:ssm:eu-west-1:886477354405:parameter/key1\nvalues:key4,arn:aws:ssm:eu-west-1:886477354405:parameter/key4\nrandomstuff"), 0644)
			input, _ := fs.Open("test.yaml.cache")
			tree := getTree()
			tree.LoadCache(input)
			Expect(len(tree.Cache)).To(Equal(2))
			Expect(tree.Cache).To(HaveKey("values:key1"))
			Expect(tree.Cache).To(HaveKey("values:key4"))
		})
	})
})
var _ = Describe("walkBranch", func() {
	Context("When the tree contains a nil key", func() {
		It("should return an error", func() {
			tree := getTree()
			branch := tree.Branches[0]
			branch[0].Key = nil
			_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (v interface{}, err error) {
				return in, err
			})
			Expect(err).To(HaveOccurred())
		})
	})
})
