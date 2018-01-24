package server

import (
	"net/http"
	"runtime"

	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/controller"
	"github.com/upsub/dispatcher/src/controller/v1"
	"github.com/upsub/dispatcher/src/dispatcher"
)

// Listen starts the http server and creates a new dispatcher
func Listen() {
	config := config.Create()
	dispatcher := dispatcher.Create(config)
	go dispatcher.Serve()

	server := &http.Server{
		ReadTimeout:  config.ConnectionTimeout,
		WriteTimeout: config.ConnectionTimeout,
		Addr:         ":" + config.Port,
	}

	http.HandleFunc("/", authenticate(config, dispatcher, controller.UpgradeHandler))
	http.HandleFunc("/v1/send", authenticate(config, dispatcher, v1.Send))
	server.ListenAndServe()
	runtime.Goexit()
}
