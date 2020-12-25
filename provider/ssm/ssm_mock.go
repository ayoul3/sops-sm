package ssm

import (
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// MockSecretValue is a dummy value for the mocks
const MockSecretValue string = "@MY_SECRET_VALUE@"

// MockClient is an AWS SSM client mock
type MockClient struct {
	ssmiface.SSMAPI
	GetParameterShouldFail bool
	ShouldBeEmpty          bool
	SecretValue            string
}

// GetSecretValue is a mock implementation of ssm method
func (m *MockClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	if m.GetParameterShouldFail {
		return nil, errors.New("GetParameter was forced to fail")
	}
	secret := MockSecretValue
	if m.SecretValue != "" {
		secret = m.SecretValue
	}
	output := new(ssm.GetParameterOutput).SetParameter(new(ssm.Parameter).SetValue(secret))
	return output, nil
}

func b64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
