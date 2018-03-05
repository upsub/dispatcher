package config

import (
	"os"
	"testing"
	"time"
)

func TestCreateDefaultConfig(t *testing.T) {
	config := Create()

	if config.Port != "4400" {
		t.Error("Port wasn't set to its default value")
	}

	if config.MaxMessageSize != 512000 {
		t.Error("MaxMessageSize should be 512000 as default")
	}

	if config.ConnectionTimeout != 30*time.Second {
		t.Error("ConnectionTimeout wasn't set to its default value")
	}

	if config.ReadTimeout != 30*time.Second {
		t.Error("ReadTimeout wasn't set to its default value")
	}

	if config.WriteTimeout != 30*time.Second {
		t.Error("WriteTimeout wasn't set to its default value")
	}

	if config.PingInterval != (9*30)/10*time.Second {
		t.Error("PingInterval wasn't set to its default value")
	}

	if config.Nats != nil {
		t.Error("Nats shouldn't be configured as default")
	}
}

func TestCreateCustomConfig(t *testing.T) {
	os.Setenv("PORT", "4401")
	os.Setenv("MAX_MESSAGE_SIZE", "10")
	os.Setenv("CONNECTION_TIMEOUT", "5")
	os.Setenv("READ_TIMEOUT", "5")
	os.Setenv("WRITE_TIMEOUT", "5")
	os.Setenv("AUTH_APP_ID", "root")
	os.Setenv("AUTH_SECRET", "secret")
	os.Setenv("AUTH_PUBLIC", "public")
	os.Setenv("AUTH_ORIGINS", "http://localhost")

	config := Create()

	if config.Port != "4401" {
		t.Error("Port wasn't set to its new value")
	}

	if config.MaxMessageSize != 10 {
		t.Error("MaxMessageSize should be set")
	}

	if config.ConnectionTimeout != 5*time.Second {
		t.Error("ConnectionTimeout wasn't set to its new value")
	}

	if config.ReadTimeout != 5*time.Second {
		t.Error("ReadTimeout wasn't set to its new value")
	}

	if config.WriteTimeout != 5*time.Second {
		t.Error("WriteTimeout wasn't set to its new value")
	}

	if config.PingInterval != (9*5*time.Second)/10 {
		t.Error("PingInterval wasn't set to its new value")
	}
}

func TestCreateInvalidConfig(t *testing.T) {
	os.Setenv("PORT", "invalid")
	os.Setenv("MAX_MESSAGE_SIZE", "invalid")
	os.Setenv("CONNECTION_TIMEOUT", "invalid")
	os.Setenv("READ_TIMEOUT", "invalid")
	os.Setenv("WRITE_TIMEOUT", "invalid")

	config := Create()

	if config.Port != "4400" {
		t.Error("Port wasn't set to its default value")
	}

	if config.MaxMessageSize != 512000 {
		t.Error("MaxMessageSize should be nil as default")
	}

	if config.ConnectionTimeout != 30*time.Second {
		t.Error("ConnectionTimeout wasn't set to its default value")
	}

	if config.ReadTimeout != 30*time.Second {
		t.Error("ReadTimeout wasn't set to its default value")
	}

	if config.WriteTimeout != 30*time.Second {
		t.Error("WriteTimeout wasn't set to its default value")
	}

	if config.PingInterval != (9*30)/10*time.Second {
		t.Error("PingInterval wasn't set to its default value")
	}
}
