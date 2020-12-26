package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ayoul3/sops-sm/provider/sm"
	"github.com/ayoul3/sops-sm/provider/ssm"
	"github.com/prometheus/common/log"
)

const SSMPatern = "arn:aws:ssm:"
const SMPatern = "arn:aws:secretsmanager:"
const defaultRegion = "eu-west-1"

type API interface {
	GetSecret(key string) (secret string, err error)
	IsSecret(key string) bool
}

type Provider struct {
	sm  API
	ssm API
}

func Init() *Provider {
	return &Provider{
		sm:  sm.NewClient(sm.NewAPI()),
		ssm: ssm.NewClient(ssm.NewAPI()),
	}
}

func (p *Provider) GetSecret(key string) (string, error) {
	region := extractRegion(key)
	if strings.Contains(key, SSMPatern) {
		p.ssm = ssm.NewClient(ssm.NewAPIForRegion(region))
		return p.ssm.GetSecret(key)
	}
	if strings.Contains(key, SMPatern) {
		p.sm = sm.NewClient(sm.NewAPIForRegion(region))
		return p.sm.GetSecret(key)
	}
	return "", fmt.Errorf("Unknown provider for secret %s", key)
}

func (p *Provider) IsSecret(key string) bool {
	return p.sm.IsSecret(key) || p.ssm.IsSecret(key)
}

func extractRegion(key string) (region string) {
	var re = regexp.MustCompile(`arn:aws:ssm:([a-z0-9-]+):\d+:`)
	match := re.FindStringSubmatch(key)
	if len(match) < 1 {
		log.Warnf("Badly formatted key %s. Could not get region.", key)
		return defaultRegion
	}
	return match[1]
}
