package config

import (
	"os"
	"testing"
)

func TestCreateNatsConfig(t *testing.T) {
	os.Setenv("NATS_HOST", "127.0.0.1")
	os.Setenv("NATS_PORT", "1000")

	nats := createNatsConfig()

	if nats.Host != "127.0.0.1" {
		t.Error("NatsConfig.Host wasn't set correctly")
	}

	if nats.Port != "1000" {
		t.Error("NatsConfig.Port wasn't set correctly")
	}
}
