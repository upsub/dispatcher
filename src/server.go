package server

import (
	"net/http"

	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/controller"
	"github.com/upsub/dispatcher/src/controller/v1"
	"github.com/upsub/dispatcher/src/dispatcher"
)

// Listen starts the http server and creates a new dispatcher
func Listen() {
	config := config.Create()
	dispatcher := dispatcher.Create()
	go dispatcher.Serve()

	http.HandleFunc("/", authenticate(config, dispatcher, controller.UpgradeHandler))
	http.HandleFunc("/v1/send", authenticate(config, dispatcher, v1.Send))
	http.ListenAndServe(":"+config.Port, nil)
}
