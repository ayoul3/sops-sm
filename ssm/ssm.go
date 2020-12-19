package myssm

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
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
	ssmAPI ssmiface.SSMAPI
}

// NewClient returns a new Client from an AWS SM client
func NewClient(api ssmiface.SSMAPI) *Client {
	return &Client{
		api,
	}
}

// NewAPI returns a new concrete AWS SSM client
func NewAPI() *ssm.SSM {
	return ssm.New(session.New())
}

// NewAPIForRegion returns a new concrete AWS SM client for a specific region
func NewAPIForRegion(region string) ssmiface.SSMAPI {
	return ssm.New(session.NewFromRegion(region))
}

// GetSecret return a Secret fetched from SM
func (c *Client) GetSecret(key string) (string, error) {
	var secret string
	res, err := c.ssmAPI.GetParameter(new(ssm.GetParameterInput).SetName(key).SetWithDecryption(true))
	if err != nil {
		return secret, err
	}
	return safeStr(res.Parameter.Value)
}

func b64Decode(input string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(input)
	return string(res), err
}

func safeStr(pointer *string) (string, error) {
	if pointer == nil {
		return "", fmt.Errorf("Null pointer in secret retrieved")
	}
	return *pointer, nil
}
