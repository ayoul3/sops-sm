package sops

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func ExtractKeyWhenJson(key, value string) (out string, err error) {
	var parsed map[string]string

	if !strings.Contains(key, "@") {
		return value, nil
	}
	keyParts := strings.Split(key, "@")
	desiredKey := keyParts[1]

	if err = json.Unmarshal([]byte(value), &parsed); err != nil {
		return "", errors.Wrap(err, "ExtractKeyWhenJson: Only simple Json structured secrets are accepted ")
	}
	for k, v := range parsed {
		if k == desiredKey {
			return v, nil
		}
	}
	return "", fmt.Errorf("ExtractKeyWhenJson: key %s not found in Json value", key)
}
