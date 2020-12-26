package sm

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/ayoul3/sops-sm/provider/session"
	"github.com/prometheus/common/log"
)

// Secret contains a SM secret details
type Secret struct {
	Key   string
	Value string
	Tags  map[string]string
}

// Client is a SM custom client
type Client struct {
	api    secretsmanageriface.SecretsManagerAPI
	region string
}

func getRegion() string {
	if os.Getenv("AWS_REGION") != "" {
		return os.Getenv("AWS_REGION")
	}
	return "eu-west-1"
}

// NewClient returns a new Client from an AWS SM client
func NewClient(api secretsmanageriface.SecretsManagerAPI, region string) *Client {
	return &Client{
		api,
		region,
	}
}

// NewAPI returns a new concrete AWS SM client
func NewAPI() (*secretsmanager.SecretsManager, string) {
	return secretsmanager.New(session.New()), getRegion()
}

// NewAPIForRegion returns a new concrete AWS SM client for a specific region
func NewAPIForRegion(region string) (secretsmanageriface.SecretsManagerAPI, string) {
	return secretsmanager.New(session.NewFromRegion(region)), region
}

// GetSecret return a Secret fetched from SM
func (c *Client) GetSecret(key string) (secret string, err error) {
	var api secretsmanageriface.SecretsManagerAPI
	if api, err = c.WithRegion(key); err != nil {
		return "", err
	}
	formattedKey := c.ExtractPath(key)
	res, err := api.GetSecretValue(new(secretsmanager.GetSecretValueInput).SetSecretId(formattedKey))
	if err != nil {
		return "", err
	}
	//secret.Key = key
	//secret.Value = *res.SecretString
	return *res.SecretString, nil
}

func (c *Client) WithRegion(key string) (secretsmanageriface.SecretsManagerAPI, error) {
	var re = regexp.MustCompile(`arn:aws:ssm:([a-z0-9-]+):\d+:`)
	match := re.FindStringSubmatch(key)
	if len(match) < 1 {
		return nil, fmt.Errorf("Badly formatted key %s. Could not get region.", key)
	}
	newRegion := match[1]
	if newRegion != c.region {
		log.Infof("Changing regions to %s", newRegion)
		newAPI, _ := NewAPIForRegion(newRegion)
		return newAPI, nil
	}
	return c.api, nil
}

func (c *Client) IsSecret(key string) bool {
	return strings.Contains(key, "arn:aws:secretsmanager")
}

func (c *Client) ExtractPath(key string) (out string) {
	var re = regexp.MustCompile(`arn:aws:secretsmanager:[a-z0-9-]+:\d+:parameter([a-zA-Z0-9/._-]+)`)
	match := re.FindStringSubmatch(key)
	if len(match) < 1 {
		log.Warnf("Badly formatted key %s", key)
		return ""
	}
	return match[1]
}
