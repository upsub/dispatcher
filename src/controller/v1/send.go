package v1

import (
	"io/ioutil"
	"net/http"

	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/dispatcher"
	"github.com/upsub/dispatcher/src/message"
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

	d.ProcessMessage(message, nil)
}
