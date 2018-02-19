package config

import "os"

// NatsConfig configuration for the NATS message bus
type NatsConfig struct {
	Host string
	Port string
}

func createNatsConfig() *NatsConfig {
	var (
		host = "localhost"
		port = "4222"
		set  = false
	)

	if value, ok := os.LookupEnv("NATS_HOST"); ok {
		set = true
		host = value
	}

	if value, ok := os.LookupEnv("NATS_PORT"); ok {
		set = true
		port = value
	}

	if !set {
		return nil
	}

	return &NatsConfig{host, port}
}
