package sm

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/ayoul3/sops-sm/provider/session"
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
func (c *Client) GetSecret(key string) (string, error) {
	formattedKey := c.ExtractPath(key)
	res, err := c.smAPI.GetSecretValue(new(secretsmanager.GetSecretValueInput).SetSecretId(formattedKey))
	if err != nil {
		return "", err
	}
	//secret.Key = key
	//secret.Value = *res.SecretString
	return *res.SecretString, nil
}

func (c *Client) IsSecret(key string) bool {
	return strings.Contains(key, "arn:aws:secretsmanager")
}

func (c *Client) ExtractPath(in string) (out string) {
	var re = regexp.MustCompile(`arn:aws:secretsmanager:[a-z0-1-]+:\d+:secret`)
	return re.ReplaceAllString(in, ``)
}
