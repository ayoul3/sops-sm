package ssm

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/ayoul3/sops-sm/provider/session"
	log "github.com/sirupsen/logrus"
)

// Secret contains a SM secret details
type Secret struct {
	Key   string
	Value string
	Tags  map[string]string
}

// Client is a SM custom client
type Client struct {
	api    ssmiface.SSMAPI
	region string
}

func getRegion() string {
	if os.Getenv("AWS_REGION") != "" {
		return os.Getenv("AWS_REGION")
	}
	return "eu-west-1"
}

// NewClient returns a new Client from an AWS SM client
func NewClient(api ssmiface.SSMAPI, region string) *Client {
	return &Client{
		api,
		region,
	}
}

// NewAPI returns a new concrete AWS SSM client
func NewAPI() (*awsssm.SSM, string) {
	return awsssm.New(session.New()), getRegion()
}

// NewAPIForRegion returns a new concrete AWS SM client for a specific region
func NewAPIForRegion(region string) (ssmiface.SSMAPI, string) {
	return awsssm.New(session.NewFromRegion(region)), region
}

// GetSecret return a Secret fetched from SM
func (c *Client) GetSecret(key string) (secret string, err error) {
	if c.api, err = c.WithRegion(key); err != nil {
		return "", err
	}
	formattedKey := c.ExtractPath(key)
	fmt.Println(formattedKey)
	res, err := c.api.GetParameter(new(awsssm.GetParameterInput).SetName(formattedKey).SetWithDecryption(true))
	if err != nil {
		return
	}
	return aws.StringValue(res.Parameter.Value), nil
}

func (c *Client) WithRegion(key string) (ssmiface.SSMAPI, error) {
	var re = regexp.MustCompile(`arn:aws:ssm:([a-z0-9-]+):\d+:`)
	match := re.FindStringSubmatch(key)
	if len(match) < 1 {
		return nil, fmt.Errorf("Badly formatted key %s. Could not get region.", key)
	}
	region := match[1]
	if region != c.region {
		log.Infof("Changing regions to %s", region)
		newAPI, _ := NewAPIForRegion(region)
		c.region = region
		return newAPI, nil
	}
	return c.api, nil
}

func (c *Client) IsSecret(key string) bool {
	return strings.Contains(key, "arn:aws:ssm:")
}

func (c *Client) ExtractPath(in string) (out string) {
	var re = regexp.MustCompile(`arn:aws:ssm:[a-z0-9-]+:\d+:parameter`)
	return re.ReplaceAllString(in, ``)
}
