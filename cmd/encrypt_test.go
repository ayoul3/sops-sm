package cmd_test

import (
	"github.com/ayoul3/sops-sm/cmd"
	"github.com/ayoul3/sops-sm/provider/ssm"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadPlainFile", func() {
	Context("When loading a plain file succeeds", func() {
		It("should return a tree", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			afero.WriteFile(handler.Fs, "test.yaml", []byte(PlainFile), 0644)
			afero.WriteFile(handler.Fs, "test.yaml.cache", []byte(cacheFile), 0644)
			loader, _ := handler.GetStore("test.yaml")
			tree, err := cmd.LoadPlainFile(handler, loader)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(tree.Branches[0])).To(Equal(1))
		})
	})
	Context("When cache file does not exist", func() {
		It("should return an error", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			afero.WriteFile(handler.Fs, "test.yaml", []byte(PlainFile), 0644)
			loader, _ := handler.GetStore("test.yaml")
			_, err := cmd.LoadPlainFile(handler, loader)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("When file does not exist", func() {
		It("should return an error", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			loader, _ := handler.GetStore("test.yaml")
			_, err := cmd.LoadPlainFile(handler, loader)
			Expect(err).To(HaveOccurred())
		})
	})
})
var _ = Describe("EncryptTree", func() {
	Context("When encrypting a tree succeeds", func() {
		It("should return bytes", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			provider := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			loader, _ := handler.GetStore("test.yaml")
			tree := getTree()

			content, err := cmd.EncryptTree(provider, loader, &tree)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]byte("name: arn:aws:ssm:eu-west-1:123456789123:parameter/someparam\n")))
		})
	})
})
var _ = Describe("DumpPlainFile", func() {
	Context("When saving file suceeds", func() {
		It("should create a file", func() {
			handler := &cmd.Handler{Fs: afero.NewMemMapFs()}
			provider := ssm.NewClient(&ssm.MockClient{SecretValue: "test"})
			loader, _ := handler.GetStore("test.yaml")
			tree := getTree()
			content, err := cmd.EncryptTree(provider, loader, &tree)
			Expect(err).ToNot(HaveOccurred())

			err = cmd.DumpPlainFile(handler, loader.GetFilePath(), content)
			Expect(err).ToNot(HaveOccurred())

			c, err := afero.ReadFile(handler.Fs, "test.yaml")
			Expect(c).To(Equal(content))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
