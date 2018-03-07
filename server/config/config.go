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
	AuthDataPath      string
	Nats              *NatsConfig
}

// Create configuration
func Create() *Config {
	var (
		port                    = "4400"
		maxMessageSize    int64 = 512 * 1000
		connectionTimeout       = 30 * time.Second
		readTimeout             = 30 * time.Second
		writeTimeout            = 30 * time.Second
	)

	if value, ok := os.LookupEnv("PORT"); ok {
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid PORT given")
		} else {
			port = value
		}
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

	return &Config{
		Port:              port,
		MaxMessageSize:    maxMessageSize,
		ConnectionTimeout: connectionTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		PingInterval:      (readTimeout * 9) / 10,
		AuthDataPath:      "auths.gob",
		Nats:              createNatsConfig(),
	}
}
