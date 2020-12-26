package sm_test

import (
	"testing"

	"github.com/ayoul3/sops-sm/provider/sm"
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
				_, err := client.GetSecret("arn:aws:secretsmanager:eu-west-1:123456789123:secret/key1")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When the keys contains an index", func() {
			It("it should return the secret", func() {
				client := sm.NewClient(&sm.MockClient{})
				secret, err := client.GetSecret("arn:aws:secretsmanager:eu-west-1:123456789123:secret/key1@index")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(sm.MockSecretValue))
			})
		})
		Context("When the call succeeds", func() {
			It("it should return the secret", func() {
				client := sm.NewClient(&sm.MockClient{})
				secret, err := client.GetSecret("arn:aws:secretsmanager:eu-west-1:123456789123:secret/key1")
				Expect(err).ToNot(HaveOccurred())
				Expect(secret).To(Equal(sm.MockSecretValue))
			})
		})
	})
})
