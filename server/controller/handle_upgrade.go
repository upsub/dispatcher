package controller

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/dispatcher"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// UpgradeHandler http request to websocket protocol
func UpgradeHandler(
	config *config.Config,
	d *dispatcher.Dispatcher,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	_, wildcard := r.Header["Origin"]

	dispatcher.CreateConnection(
		r.Header.Get("Sec-Websocket-Key"),
		r.Header.Get("upsub-app-id"),
		r.Header.Get("upsub-connection-name"),
		conn,
		config,
		d,
		map[string]bool{
			"wildcard": !wildcard,
		},
	)
}
