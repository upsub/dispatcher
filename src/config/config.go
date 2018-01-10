package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the main configuration struct
type Config struct {
	Port           string
	MaxMessageSize int64
	Timeout        time.Duration
	WriteTimeout   time.Duration
	PingInterval   time.Duration
	Auths          map[string]*AuthConfig
	Nats           *NatsConfig
}

// NatsConfig configuration for the NATS message bus
type NatsConfig struct {
	Host string
	Port string
}

// AuthConfig handles authentication configuration
type AuthConfig struct {
	AppID   string
	Secret  string
	Public  string
	Origins []string
}

func createNatsConfig() *NatsConfig {
	var (
		host = "localhost"
		port = "6379"
	)

	if value, ok := os.LookupEnv("NATS_HOST"); ok {
		host = value
	}

	if value, ok := os.LookupEnv("NATS_PORT"); ok {
		port = value
	}

	return &NatsConfig{host, port}
}

func createAuthConfig() *AuthConfig {
	config := &AuthConfig{}

	if value, ok := os.LookupEnv("AUTH_APP_ID"); ok {
		config.AppID = value
	}

	if value, ok := os.LookupEnv("AUTH_SECRET"); ok {
		config.Secret = value
	}

	if value, ok := os.LookupEnv("AUTH_PUBLIC"); ok {
		config.Public = value
	}

	if value, ok := os.LookupEnv("AUTH_ORIGINS"); ok {
		config.Origins = strings.Split(value, ",")
	}

	return config
}

func Create() *Config {
	var (
		port           = "5000"
		maxMessageSize int64
		timeout        = 10 * time.Second
		writeTimeout   = 10 * time.Second
		pingInterval   = (timeout * 9) / 10
		auths          = map[string]*AuthConfig{}
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

	if value, ok := os.LookupEnv("TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid TIMEOUT given")
		} else {
			timeout = time.Duration(converted) * time.Second
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
		Port:           port,
		MaxMessageSize: maxMessageSize,
		Timeout:        timeout,
		WriteTimeout:   writeTimeout,
		PingInterval:   pingInterval,
		Auths:          auths,
		Nats:           createNatsConfig(),
	}
}
