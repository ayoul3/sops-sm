package provider

import (
	"fmt"
	"regexp"

	"github.com/ayoul3/sops-sm/provider/sm"
	"github.com/ayoul3/sops-sm/provider/ssm"
	"github.com/pkg/errors"
)

const defaultRegion = "eu-west-1"

type API interface {
	GetSecret(key string) (secret string, err error)
	IsSecret(key string) bool
	WithRegion(region string)
}

type Provider struct {
	Apis map[string]API
}

func Init() *Provider {
	p := &Provider{}
	p.Apis = map[string]API{
		ssm.SSMPatern: ssm.NewClient(ssm.NewAPI()),
		sm.SMPatern:   sm.NewClient(sm.NewAPI()),
	}
	return p
}
func (p *Provider) WithRegion(region string) {
}

func (p *Provider) GetSecret(key string) (secret string, err error) {
	var pattern, region string
	if pattern, region, err = ExtractRegionPattern(key); err != nil {
		return "", errors.Wrapf(err, "GetSecret ")
	}
	if client, ok := p.Apis[pattern]; ok {
		client.WithRegion(region)
		return client.GetSecret(key)
	}
	return "", fmt.Errorf("Unknown provider for secret %s", key)
}

func (p *Provider) IsSecret(key string) bool {
	for _, client := range p.Apis {
		if client.IsSecret(key) {
			return true
		}
	}
	return false
}

func ExtractRegionPattern(key string) (pattern, region string, err error) {
	var re = regexp.MustCompile(`(arn:aws:(?:ssm|secretsmanager):)([a-z0-9-]+):\d+:`)
	match := re.FindStringSubmatch(key)
	if len(match) < 2 {
		return "", "", fmt.Errorf("ExtractRegionPattern: Badly formatted key %s. Could not get pattern and region.", key)
	}
	return match[1], match[2], nil
}
