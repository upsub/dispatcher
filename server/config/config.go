package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

// Config is the main configuration struct
type Config struct {
	Port              int64
	MaxMessageSize    int64         `toml:"max-message-size"`
	ConnectionTimeout time.Duration `toml:"connection-timeout"`
	ReadTimeout       time.Duration `toml:"read-timeout"`
	WriteTimeout      time.Duration `toml:"write-timeout"`
	PingInterval      time.Duration
	AuthDataPath      string
	ConfigPath        string
	Nats              *NatsConfig
}

// Create configuration
func Create() *Config {
	defaultTimeout := 30 * time.Second
	config := &Config{
		Port:              4400,
		MaxMessageSize:    512 * 1000,
		ConnectionTimeout: defaultTimeout,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		PingInterval:      (defaultTimeout * 9) / 10,
		AuthDataPath:      "auths.gob",
		ConfigPath:        "config.toml",
		Nats:              nil,
	}

	config = load(config)

	if value, ok := os.LookupEnv("PORT"); ok {
		if port, err := strconv.ParseInt(value, 10, 64); err == nil {
			config.Port = port
		} else {
			log.Print("Invalid PORT given")
		}
	}

	if value, ok := os.LookupEnv("MAX_MESSAGE_SIZE"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid MAX_MESSAGE_SIZE given")
		} else {
			config.MaxMessageSize = converted
		}
	}

	if value, ok := os.LookupEnv("CONNECTION_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid CONNECTION_TIMEOUT given")
		} else {
			config.ConnectionTimeout = time.Duration(converted) * time.Second
		}
	}

	if value, ok := os.LookupEnv("READ_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid READ_TIMEOUT given")
		} else {
			config.ReadTimeout = time.Duration(converted) * time.Second
		}
	}

	if value, ok := os.LookupEnv("WRITE_TIMEOUT"); ok {
		if converted, err := strconv.ParseInt(value, 10, 64); err != nil {
			log.Print("Invalid WRITE_TIMEOUT given")
		} else {
			config.WriteTimeout = time.Duration(converted) * time.Second
		}
	}

	config.PingInterval = (config.ReadTimeout * 9) / 10
	config.Nats = createNatsConfig()
	return config
}

func load(config *Config) *Config {
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		return config
	}

	buffer, err := ioutil.ReadFile(config.ConfigPath)

	if err != nil {
		log.Print("[CONFIG ERROR] ", err)
		return config
	}

	if _, err := toml.Decode(string(buffer), config); err != nil {
		log.Print("[CONFIG ERROR] ", err)
	}

	config.ConnectionTimeout = config.ConnectionTimeout * time.Second
	config.WriteTimeout = config.WriteTimeout * time.Second
	config.ReadTimeout = config.ReadTimeout * time.Second
	config.PingInterval = (config.ReadTimeout * 9) / 10
	return config
}
