package server

import (
	"net/http"

	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/dispatcher"
)

// Listen starts the dispatcher
func Listen() {
	config := config.Create()
	dispatcher := dispatcher.Create()
	go dispatcher.Serve()

	http.HandleFunc("/", authenticate(config, dispatcher, dispatcher.UpgradeHandler))
	http.ListenAndServe(":"+config.Port, nil)
}
