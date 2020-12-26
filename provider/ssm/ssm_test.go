package ssm_test

import (
	"testing"

	"github.com/ayoul3/sops-sm/provider/ssm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestAWS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "AWS/SSM", []Reporter{reporters.NewJUnitReporter("test_report-aws-ssm.xml")})
}

var _ = Describe("SSM", func() {
	Describe("GetSecret", func() {
		Context("When the client fails", func() {
			It("should return an error", func() {
				client := ssm.NewClient(&ssm.MockClient{GetParameterShouldFail: true})
				_, err := client.GetSecret("arn:aws:ssm:eu-west-1:886477354405:parameter/key1")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When the keys contains an index", func() {
			It("it should return the secret", func() {
				client := ssm.NewClient(&ssm.MockClient{})
				secret, err := client.GetSecret("arn:aws:ssm:eu-west-1:886477354405:parameter/key1@index")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(ssm.MockSecretValue))
			})
		})
		Context("When the call succeeds", func() {
			It("it should return the secret", func() {
				client := ssm.NewClient(&ssm.MockClient{})
				secret, err := client.GetSecret("arn:aws:ssm:eu-west-1:886477354405:parameter/key1")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(ssm.MockSecretValue))
			})
		})
	})
})
