package myssm_test

import (
	"testing"

	myssm "github.com/ayoul3/sops-sm/ssm"
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
				client := myssm.NewClient(&myssm.MockClient{GetParameterShouldFail: true})
				_, err := client.GetSecret("test")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When the call succeeds", func() {
			It("it should return the secret", func() {
				client := myssm.NewClient(&myssm.MockClient{})
				secret, err := client.GetSecret("test")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(myssm.MockSecretValue))
			})
		})
	})
})
