package dispatcher

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	port           string
	maxMessageSize int64
	timeout        time.Duration
	writeTimeout   time.Duration
	pingInterval   time.Duration
	auths          map[string]*authConfig
	redis          *natsConfig
}

type natsConfig struct {
	host string
	port string
}

type authConfig struct {
	appID   string
	secret  string
	public  string
	origins []string
}

func createNatsConfig() *natsConfig {
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

	return &natsConfig{host, port}
}

func createAuthConfig() *authConfig {
	config := &authConfig{}

	if value, ok := os.LookupEnv("AUTH_APP_ID"); ok {
		config.appID = value
	}

	if value, ok := os.LookupEnv("AUTH_SECRET"); ok {
		config.secret = value
	}

	if value, ok := os.LookupEnv("AUTH_PUBLIC"); ok {
		config.public = value
	}

	if value, ok := os.LookupEnv("AUTH_ORIGINS"); ok {
		config.origins = strings.Split(value, ",")
	}

	return config
}

func createConfig() *config {
	var (
		port           = "5000"
		maxMessageSize int64
		timeout        = 10 * time.Second
		writeTimeout   = 10 * time.Second
		pingInterval   = (timeout * 9) / 10
		auths          = map[string]*authConfig{}
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

	if auth := createAuthConfig(); auth.appID != "" {
		auths[auth.appID] = auth
	}

	return &config{
		port:           port,
		maxMessageSize: maxMessageSize,
		timeout:        timeout,
		writeTimeout:   writeTimeout,
		pingInterval:   pingInterval,
		auths:          auths,
		redis:          createNatsConfig(),
	}
}
