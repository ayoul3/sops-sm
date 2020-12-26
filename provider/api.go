package provider

import (
	"fmt"
	"strings"

	"github.com/ayoul3/sops-sm/provider/sm"
	"github.com/ayoul3/sops-sm/provider/ssm"
)

const SSMPatern = "arn:aws:ssm:"
const SMPatern = "arn:aws:secretsmanager:"

type API interface {
	GetSecret(key string) (string, error)
	IsSecret(key string) bool
}

type Provider struct {
	ssm API
	sm  API
}

func Init() *Provider {
	return &Provider{
		ssm: ssm.NewClient(ssm.NewAPI()),
		sm:  sm.NewClient(sm.NewAPI()),
	}
}

func (p *Provider) GetSecret(key string) (string, error) {
	if strings.Contains(key, SSMPatern) {
		return p.ssm.GetSecret(key)
	}
	if strings.Contains(key, SMPatern) {
		return p.ssm.GetSecret(key)
	}
	return "", fmt.Errorf("Unknown provider for secret %s", key)
}

func (p *Provider) IsSecret(key string) bool {
	return p.sm.IsSecret(key) || p.ssm.IsSecret(key)
}
