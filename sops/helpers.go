package sops

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ayoul3/sops-sm/provider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func PrepareAsync(tree *Tree, provider provider.API, numThreads int) {
	InitWorkers(numThreads)
	go RunWorkers(provider)
	go CacheAsyncSecret(tree)
}

func ExtractKeyWhenJson(key, value string) (out string, err error) {
	var parsed map[string]string

	if !strings.Contains(key, "#") {
		return value, nil
	}
	keyParts := strings.Split(key, "#")
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

var WalkerAsyncFetchSecret = func(branch TreeBranch, provider provider.API) error {
	_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (v interface{}, err error) {
		var ok bool

		pathString := strings.Join(path, ":")
		log.Infof("Walking path %s ", pathString)
		if v, ok = in.(string); !ok {
			return in, nil
		}
		if provider.IsSecret(v.(string)) {
			log.Infof("sending secret for async processing %s", in.(string))
			MsgChan <- WorkerSecret{Key: in.(string)}
		}
		return v, nil
	})
	return err
}

var WalkerSyncFetchSecret = func(tree *Tree, branch TreeBranch, provider provider.API) error {
	_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (v interface{}, err error) {
		var cached, secretValue string
		var ok, found bool

		pathString := strings.Join(path, ":")
		log.Infof("Walking path %s ", pathString)

		if v, ok = in.(string); !ok {
			return in, nil
		}
		if cached, found = tree.IsCached(v.(string)); found {
			log.Infof("Found secret in cache %s", v)
			tree.CacheSecretValue(v.(string), cached, pathString) // update cache path
			return ExtractKeyWhenJson(v.(string), cached)
		}
		if provider.IsSecret(v.(string)) {
			log.Infof("Fetching secret %s ", v)
			if secretValue, err = provider.GetSecret(v.(string)); err != nil {
				return nil, err
			}
			tree.CacheSecretValue(v.(string), secretValue, pathString)
			return ExtractKeyWhenJson(v.(string), secretValue)
		}
		return v, nil

	})
	return err
}

var WalkerEncryptSecret = func(tree *Tree, branch TreeBranch, provider provider.API) error {
	_, err := branch.walkBranch(branch, make([]string, 0), func(in interface{}, path []string) (v interface{}, err error) {
		var cached CachedSecret
		var ok, found bool

		pathString := strings.Join(path, ":")
		log.Infof("Walking path %s ", pathString)

		if v, ok = in.(string); !ok {
			return in, nil
		}
		if cached, found = tree.Cache[pathString]; found {
			log.Infof("Found secret in cache %s", v)
			return cached.Value, nil
		}
		return v, nil

	})
	return err
}
