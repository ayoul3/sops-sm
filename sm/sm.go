package sm

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/ayoul3/sops-sm/session"
)

// Secret contains a SM secret details
type Secret struct {
	Key   string
	Value string
	Tags  map[string]string
}

// Client is a SM custom client
type Client struct {
	smAPI secretsmanageriface.SecretsManagerAPI
}

// NewClient returns a new Client from an AWS SM client
func NewClient(api secretsmanageriface.SecretsManagerAPI) *Client {
	return &Client{
		api,
	}
}

// NewAPI returns a new concrete AWS SM client
func NewAPI() *secretsmanager.SecretsManager {
	return secretsmanager.New(session.New())
}

// NewAPIForRegion returns a new concrete AWS SM client for a specific region
func NewAPIForRegion(region string) secretsmanageriface.SecretsManagerAPI {
	return secretsmanager.New(session.NewFromRegion(region))
}

// GetSecret return a Secret fetched from SM
func (c *Client) GetSecret(key string) (Secret, error) {
	var secret Secret
	res, err := c.smAPI.GetSecretValue(new(secretsmanager.GetSecretValueInput).SetSecretId(key))
	if err != nil {
		return secret, err
	}
	secret.Key = key
	secret.Value = *res.SecretString
	return secret, nil
}

// GetSecretWithTags return a Secret fetched from SM with its tags
func (c *Client) GetSecretWithTags(key string) (Secret, error) {
	secret, err := c.GetSecret(key)
	if err != nil {
		return secret, err
	}
	res, err := c.smAPI.DescribeSecret(new(secretsmanager.DescribeSecretInput).SetSecretId(key))
	if err != nil {
		return secret, err
	}
	secret.Tags = make(map[string]string)
	for _, tag := range res.Tags {
		secret.Tags[*tag.Key] = *tag.Value
	}
	return secret, nil
}

func b64Decode(input string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(input)
	return string(res), err
}
