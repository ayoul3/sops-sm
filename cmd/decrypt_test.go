package cmd_test

import (
	"github.com/ayoul3/sops-sm/cmd"
	"github.com/ayoul3/sops-sm/provider/ssm"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadEncryptedFile", func() {
	Context("When loading an encrypted file succeeds", func() {
		It("should return a tree", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			afero.WriteFile(handler.Fs, "test.yaml", []byte(EncryptedFile), 0644)
			loader, _ := handler.GetStore("test.yaml")
			tree, err := cmd.LoadEncryptedFile(handler, loader)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(tree.Branches[0])).To(Equal(1))
		})
	})
	Context("When file does not exist", func() {
		It("should return an error", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			loader, _ := handler.GetStore("test.yaml")
			_, err := cmd.LoadEncryptedFile(handler, loader)
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = Describe("DecryptTree", func() {
	Context("When decrypting a tree succeeds", func() {
		It("should return bytes", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			provider := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			loader, _ := handler.GetStore("test.yaml")
			tree := getTree()

			content, err := cmd.DecryptTree(provider, loader, &tree, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]byte("name: test\n")))
		})
	})
	Context("When decrypting fails", func() {
		It("should returnan error", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			provider := ssm.NewClient(&ssm.MockClient{GetParameterShouldFail: true})
			loader, _ := handler.GetStore("test.yaml")
			tree := getTree()

			_, err := cmd.DecryptTree(provider, loader, &tree, 1)
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = Describe("DumpDecryptedTree", func() {
	Context("When saving file suceeds", func() {
		It("should create two files", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			provider := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			loader, _ := handler.GetStore("test.yaml")
			tree := getTree()
			content, err := cmd.DecryptTree(provider, loader, &tree, 1)
			Expect(err).ToNot(HaveOccurred())

			err = cmd.DumpDecryptedTree(handler, loader.GetFilePath(), loader.GetCachePath(), content, tree.GetCache())
			Expect(err).ToNot(HaveOccurred())

			c, err := afero.ReadFile(handler.Fs, "test.yaml")
			Expect(c).To(Equal(content))
			Expect(err).ToNot(HaveOccurred())

			c, err = afero.ReadFile(handler.Fs, "test.yaml.cache")
			Expect(err).ToNot(HaveOccurred())
			Expect(c).To(Equal([]byte("name,arn:aws:ssm:eu-west-1:123456789123:parameter/someparam\n")))
		})
	})
})
