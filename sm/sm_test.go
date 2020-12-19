package sm_test

import (
	"testing"

	"github.com/ayoul3/sops-sm/sm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestAWS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "AWS/SM", []Reporter{reporters.NewJUnitReporter("test_report-aws-sm.xml")})
}

var _ = Describe("SM", func() {
	Describe("GetSecret", func() {
		Context("When the client fails", func() {
			It("should return an error", func() {
				client := sm.NewClient(&sm.MockClient{GetSecretShouldFail: true})
				_, err := client.GetSecret("test")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When the call succeeds", func() {
			It("it should return the secret", func() {
				client := sm.NewClient(&sm.MockClient{})
				secret, err := client.GetSecret("test")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret.Value).To(Equal(sm.MockSecretValue))
				Expect(secret.Key).To(Equal("test"))
			})
		})
	})
	Describe("GetSecretWithTags", func() {
		Context("When the client fails", func() {
			Context("When the failure occurs on the GetSecret", func() {
				It("should return an error", func() {
					client := sm.NewClient(&sm.MockClient{GetSecretShouldFail: true})
					_, err := client.GetSecretWithTags("test")
					Expect(err).To(HaveOccurred())
				})
			})
			Context("When the failure occurs on the DescribeSecret", func() {
				It("should return an error", func() {
					client := sm.NewClient(&sm.MockClient{DescribeSecretShouldFail: true})
					_, err := client.GetSecretWithTags("test")
					Expect(err).To(HaveOccurred())
				})
			})
		})
		Context("When the client returns the secret details", func() {
			Context("When the secret has no tags", func() {
				It("should return no tags", func() {
					client := sm.NewClient(&sm.MockClient{ShouldBeEmpty: true})
					secret, err := client.GetSecretWithTags("test")
					Expect(err).ToNot(HaveOccurred())
					Expect(secret.Value).To(Equal(sm.MockSecretValue))
					Expect(secret.Key).To(Equal("test"))
					Expect(secret.Tags).To(BeEmpty())
				})
			})
			Context("When the secret has tags", func() {
				It("should return tags", func() {
					client := sm.NewClient(&sm.MockClient{})
					secret, err := client.GetSecretWithTags("test")
					Expect(err).ToNot(HaveOccurred())
					Expect(secret.Value).To(Equal(sm.MockSecretValue))
					Expect(secret.Key).To(Equal("test"))
					Expect(secret.Tags).To(Equal(map[string]string{"MY_KEY": "MY_VALUE"}))
				})
			})
		})
	})
})
