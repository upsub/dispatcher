package controller

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/dispatcher"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	dispatcher.CreateConnection(
		r.Header.Get("Sec-Websocket-Key"),
		r.Header.Get("upsub-app-id"),
		r.Header.Get("upsub-connection-name"),
		conn,
		config,
		d,
	)
}
