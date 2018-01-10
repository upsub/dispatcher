package dispatcher

import (
	"log"
	"net/http"

	"github.com/upsub/dispatcher/src/util"
)

func validateAppID(config *config, appID string) bool {
	for id := range config.auths {
		if id == appID {
			return true
		}
	}

	return false
}

func validateSecretKey(config *config, secret string) bool {
	for _, auth := range config.auths {
		if auth.secret == secret {
			return true
		}
	}

	return false
}

func validatePublicKey(config *config, public string, origin string) bool {
	for _, auth := range config.auths {
		if auth.public == public && util.Contains(auth.origins, origin) {
			return true
		}
	}

	return false
}

type controller func(*config, *dispatcher, http.ResponseWriter, *http.Request)

func authenticate(c *config, d *dispatcher, next controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if query := r.URL.Query(); len(query) > 0 {
			r.Header.Set("upsub-app-id", query.Get("upsub-app-id"))
			r.Header.Set("upsub-public", query.Get("upsub-public"))
		}

		if len(c.auths) == 0 {
			next(c, d, w, r)
			return
		}

		if ok := validateAppID(c, r.Header.Get("upsub-app-id")); !ok {
			log.Print("Invalid APP ID")
			http.Error(w, "Invalid APP ID", 401)
			return
		}

		if ok := validateSecretKey(c, r.Header.Get("upsub-secret")); !ok && r.Header.Get("origin") == "" {
			log.Print("Invalid Secret key")
			http.Error(w, "Invalid Secret key", 401)
			return
		}

		if ok := validatePublicKey(c, r.Header.Get("upsub-public"), r.Header.Get("origin")); !ok {
			log.Print("Invalid public key or origin")
			http.Error(w, "Invalid public key or origin", 401)
			return
		}

		next(c, d, w, r)
	}
}
