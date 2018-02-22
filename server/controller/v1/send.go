package v1

import (
	"io/ioutil"
	"net/http"

	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/dispatcher"
	"github.com/upsub/dispatcher/server/message"
)

// Send is dispatching events to upsub clients
func Send(
	config *config.Config,
	d *dispatcher.Dispatcher,
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Body == nil {
		http.Error(w, "Empty request body", 400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed reading request body", 400)
		return
	}

	message, err := message.Decode(body)

	if err != nil {
		http.Error(w, "Invalid request body", 400)
		return
	}

	d.ProcessMessage(message, dispatcher.CreateConnection(
		r.Header.Get("Sec-Websocket-Key"),
		r.Header.Get("upsub-app-id"),
		r.Header.Get("upsub-connection-name"),
		nil,
		config,
		d,
		map[string]bool{
			"wildcard": false,
		},
	))
}
