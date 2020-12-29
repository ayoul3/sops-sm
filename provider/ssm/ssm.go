package ssm

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/ayoul3/sops-sm/provider/session"
	log "github.com/sirupsen/logrus"
)

const SSMPatern = "arn:aws:ssm:"

// Client is a SM custom client
type Client struct {
	api ssmiface.SSMAPI
}

// NewClient returns a new Client from an AWS SM client
func NewClient(api ssmiface.SSMAPI) *Client {
	return &Client{
		api,
	}
}

// NewAPI returns a new concrete AWS SSM client
func NewAPI() *awsssm.SSM {
	return awsssm.New(session.New())
}

// NewAPIForRegion returns a new concrete AWS SM client for a specific region
func NewAPIForRegion(region string) ssmiface.SSMAPI {
	return awsssm.New(session.NewFromRegion(region))
}

// Overrides region
func (c *Client) WithRegion(region string) {
	c.api = awsssm.New(session.NewFromRegion(region))
}

// GetSecret return a Secret fetched from SSM
func (c *Client) GetSecret(key string) (secret string, err error) {
	formattedKey := c.ExtractPath(key)
	res, err := c.api.GetParameter(new(awsssm.GetParameterInput).SetName(formattedKey).SetWithDecryption(true))
	if err != nil {
		return
	}
	return aws.StringValue(res.Parameter.Value), nil
}

func (c *Client) IsSecret(key string) bool {
	return strings.Contains(key, "arn:aws:ssm:")
}

func (c *Client) ExtractPath(key string) (out string) {
	var re = regexp.MustCompile(`arn:aws:ssm:[a-z0-9-]+:\d+:parameter([a-zA-Z0-9/._-]+)`)
	match := re.FindStringSubmatch(key)
	if len(match) < 2 {
		log.Warnf("Badly formatted key %s", key)
		return key
	}
	return match[1]
}
