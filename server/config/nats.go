package config

import (
	"log"
	"os"
	"strconv"
)

// NatsConfig configuration for the NATS message bus
type NatsConfig struct {
	Host string
	Port int64
}

func createNatsConfig() *NatsConfig {
	var (
		host = "localhost"
		port = int64(4222)
		set  = false
	)

	if value, ok := os.LookupEnv("NATS_HOST"); ok {
		set = true
		host = value
	}

	if value, ok := os.LookupEnv("NATS_PORT"); ok {
		set = true
		if num, err := strconv.ParseInt(value, 10, 64); err == nil {
			port = num
		} else {
			log.Print("Invalid NATS_PORT given")
		}
	}

	if !set {
		return nil
	}

	return &NatsConfig{host, port}
}
