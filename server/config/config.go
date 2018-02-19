package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config is the main configuration struct
type Config struct {
	Port              string
	MaxMessageSize    int64
	ConnectionTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	PingInterval      time.Duration
	Auths             map[string]*AuthConfig
	Nats              *NatsConfig
}

// Create configuration
func Create() *Config {
	var (
		port              = "4400"
		maxMessageSize    int64
		connectionTimeout = 10 * time.Second
		readTimeout       = 10 * time.Second
		writeTimeout      = 10 * time.Second
		pingInterval      = (readTimeout * 9) / 10
		auths             = map[string]*AuthConfig{}
	)

	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	if value, ok := os.LookupEnv("MAX_MESSAGE_SIZE"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid MAX_MESSAGE_SIZE given")
		} else {
			maxMessageSize = converted
		}
	}

	if value, ok := os.LookupEnv("CONNECTION_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid CONNECTION_TIMEOUT given")
		} else {
			connectionTimeout = time.Duration(converted) * time.Second
		}
	}

	if value, ok := os.LookupEnv("READ_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid READ_TIMEOUT given")
		} else {
			readTimeout = time.Duration(converted) * time.Second
		}
	}

	if value, ok := os.LookupEnv("WRITE_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid WRITE_TIMEOUT given")
		} else {
			writeTimeout = time.Duration(converted) * time.Second
		}
	}

	if auth := createAuthConfig(); auth.AppID != "" {
		auths[auth.AppID] = auth
	}

	return &Config{
		Port:              port,
		MaxMessageSize:    maxMessageSize,
		ConnectionTimeout: connectionTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		PingInterval:      pingInterval,
		Auths:             auths,
		Nats:              createNatsConfig(),
	}
}