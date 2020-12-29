package sops

import (
	"time"

	"github.com/ayoul3/sops-sm/provider"
	log "github.com/sirupsen/logrus"
)

type WorkerSecret struct {
	Path  string
	Key   string
	Value string
}

var MsgChan chan WorkerSecret
var ReportChan chan WorkerSecret
var WorkerDone chan bool

func InitWorkers(numThreads int) {
	MsgChan = make(chan WorkerSecret, numThreads)
	ReportChan = make(chan WorkerSecret, numThreads)
	WorkerDone = make(chan bool)
}

func RunWorkers(provider provider.API) {
	for {
		msg := <-MsgChan
		go func(msg WorkerSecret) {
			var err error
			log.Infof("async - decrypting secret %s", msg.Key)
			if msg.Value, err = provider.GetSecret(msg.Key); err != nil {
				log.Warnf("RunWorkers: Error fetching secret %s: %s", msg.Key, err)
				return
			}
			ReportChan <- msg
		}(msg)
	}
}

func CacheAsyncSecret(tree *Tree) {
loop:
	for {
		select {
		case msg := <-ReportChan:
			log.Infof("async - received secret for storage %s", msg.Key)
			tree.CacheSecretValue(msg.Key, msg.Value, msg.Path)
		case <-time.After(500 * time.Millisecond):
			log.Info("Finished waiting for secrets")
			break loop
		}
	}
	WorkerDone <- true
}
