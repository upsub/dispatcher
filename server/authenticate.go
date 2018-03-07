package server

import (
	"log"
	"net/http"

	"github.com/upsub/dispatcher/server/auth"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/dispatcher"
	"github.com/upsub/dispatcher/server/util"
)

type handler func(*config.Config, *dispatcher.Dispatcher, http.ResponseWriter, *http.Request)

func parseQueryParams(r *http.Request) {
	allowedQueryParams := []string{
		"upsub-app-id",
		"upsub-secret",
		"upsub-public",
		"upsub-connection-name",
	}

	if query := r.URL.Query(); len(query) > 0 {
		for _, key := range allowedQueryParams {
			r.Header.Set(key, query.Get(key))
		}
	}
}

func reject(w http.ResponseWriter) {
	log.Print("[AUTH] Invalid authentication credentials")
	http.Error(w, "Invalid authentication credentials", 401)
}

func authenticate(
	c *config.Config,
	d *dispatcher.Dispatcher,
	store *auth.Store,
	next handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parseQueryParams(r)

		if store.Length() == 0 {
			next(c, d, w, r)
			return
		}

		auth := store.Find(r.Header.Get("upsub-app-id"))

		if auth == nil {
			reject(w)
			return
		}

		origin := r.Header.Get("Origin")

		if len(origin) == 0 && auth.Secret == r.Header.Get("upsub-secret") {
			next(c, d, w, r)
			return
		}

		if len(origin) > 0 && auth.Public == r.Header.Get("upsub-public") && util.Contains(auth.Origins, origin) {
			next(c, d, w, r)
			return
		}

		reject(w)
	}
}
