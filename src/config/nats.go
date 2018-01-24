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
	)

	if value, ok := os.LookupEnv("NATS_HOST"); ok {
		host = value
	}

	if value, ok := os.LookupEnv("NATS_PORT"); ok {
		port = value
	}

	return &NatsConfig{host, port}
}
