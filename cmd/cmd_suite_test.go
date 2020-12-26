package cmd_test

import (
	"testing"

	"github.com/ayoul3/sops-sm/sops"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "cmd", []Reporter{reporters.NewJUnitReporter("test_report-cmd.xml")})
}

func getTree() sops.Tree {
	return sops.Tree{
		Cache: make(map[string]sops.CachedSecret, 0),
		Branches: sops.TreeBranches{
			sops.TreeBranch{
				sops.TreeItem{
					Key:   "name",
					Value: "arn:aws:ssm:eu-west-1:123456789123:parameter/someparam",
				},
			},
		},
	}
}

var PlainFile = []byte(`---
name: test`)
var EncryptedFile = []byte(`---
name: arn:aws:ssm:eu-west-1:123456789123:parameter/someparam`)
var cacheFile = []byte(`name,arn:aws:ssm:eu-west-1:123456789123:parameter/someparam`)
