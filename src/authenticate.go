package server

import (
	"log"
	"net/http"

	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/dispatcher"
	"github.com/upsub/dispatcher/src/util"
)

func validateAppID(config *config.Config, appID string) bool {
	for id := range config.Auths {
		if id == appID {
			return true
		}
	}

	return false
}

func validateSecretKey(config *config.Config, secret string) bool {
	for _, auth := range config.Auths {
		if auth.Secret == secret {
			return true
		}
	}

	return false
}

func validatePublicKey(config *config.Config, public string, origin string) bool {
	for _, auth := range config.Auths {
		if auth.Public == public && util.Contains(auth.Origins, origin) {
			return true
		}
	}

	return false
}

type handler func(*config.Config, *dispatcher.Dispatcher, http.ResponseWriter, *http.Request)

func authenticate(c *config.Config, d *dispatcher.Dispatcher, next handler) http.HandlerFunc {
	allowedQueryParams := []string{
		"upsub-app-id",
		"upsub-secret",
		"upsub-public",
		"upsub-connection-name",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if query := r.URL.Query(); len(query) > 0 {
			for _, key := range allowedQueryParams {
				r.Header.Set(key, query.Get(key))
			}
		}

		if len(c.Auths) == 0 {
			next(c, d, w, r)
			return
		}

		if ok := validateAppID(c, r.Header.Get("upsub-app-id")); !ok {
			log.Print("Invalid APP ID")
			http.Error(w, "Invalid APP ID", 401)
			return
		}

		_, origin := r.Header["Origin"]
		if !validateSecretKey(c, r.Header.Get("upsub-secret")) && !origin {
			log.Print("Invalid Secret key")
			http.Error(w, "Invalid Secret key", 401)
			return
		}

		if origin && !validatePublicKey(c, r.Header.Get("upsub-public"), r.Header.Get("origin")) {
			log.Print("Invalid public key or origin")
			http.Error(w, "Invalid public key or origin", 401)
			return
		}

		next(c, d, w, r)
	}
}
