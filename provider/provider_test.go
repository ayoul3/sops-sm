package provider_test

import (
	"testing"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/ayoul3/sops-sm/provider/sm"
	"github.com/ayoul3/sops-sm/provider/ssm"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestAWS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "provider", []Reporter{reporters.NewJUnitReporter("test_report-provider.xml")})
}

var _ = Describe("provider", func() {
	Describe("GetSecret", func() {
		Context("When the pattern is not found", func() {
			provider := provider.Provider{
				Apis: map[string]provider.API{
					ssm.SSMPatern: ssm.NewClient(&ssm.MockClient{}),
				},
			}
			It("should return an error about missing region", func() {
				_, err := provider.GetSecret("whateverkey")
				Expect(err).To(HaveOccurred())
			})
			It("should return an error about missing pattern", func() {
				_, err := provider.GetSecret("arn:aws:secretsmanager:eu-west-1:123456789123:secret/key1@index")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("IsSecret", func() {
		Context("when it's a secret", func() {
			provider := provider.Provider{
				Apis: map[string]provider.API{
					ssm.SSMPatern: ssm.NewClient(&ssm.MockClient{}),
					sm.SMPatern:   sm.NewClient(&sm.MockClient{}),
				},
			}
			It("should return true for SSM", func() {
				check := provider.IsSecret("arn:aws:ssm:eu-west-1:123456789123:parameter/key1@index")
				Expect(check).To(BeTrue())
			})
			It("should return true for SM", func() {
				check := provider.IsSecret("arn:aws:secretsmanager:eu-west-1:123456789123:secret/key1@index")
				Expect(check).To(BeTrue())
			})
		})
		Context("when it's not a secret", func() {
			provider := provider.Provider{
				Apis: map[string]provider.API{
					ssm.SSMPatern: ssm.NewClient(&ssm.MockClient{}),
					sm.SMPatern:   sm.NewClient(&sm.MockClient{}),
				},
			}
			It("should return false", func() {
				check := provider.IsSecret("whatever")
				Expect(check).ToNot(BeTrue())
			})
		})
	})
})
