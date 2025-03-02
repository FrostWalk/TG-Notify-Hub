package healtcheck

import (
	"github.com/etkecc/go-healthchecks/v2"
	"strings"
	"time"
)

var (
	client *healthchecks.Client = nil
)

func EnableCheck(uuid string, interval int) {
	client = healthchecks.New(healthchecks.WithCheckUUID(uuid))
	go client.Auto(time.Duration(interval) * time.Second)
}

func CloseConnection() {
	if client != nil {
		client.Shutdown()
	}
}

func SignalError(err error) {
	if client != nil {
		client.Fail(strings.NewReader(err.Error()))
	}
}
