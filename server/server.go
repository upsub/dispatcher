package server

import (
	"net/http"
	"runtime"

	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/controller"
	"github.com/upsub/dispatcher/server/controller/v1"
	"github.com/upsub/dispatcher/server/dispatcher"
)

// Listen starts the http server and creates a new dispatcher
func Listen() {
	conf := config.Create()
	dispatcher := dispatcher.Create(conf)
	go dispatcher.Serve()

	server := &http.Server{
		ReadTimeout:  conf.ConnectionTimeout,
		WriteTimeout: conf.ConnectionTimeout,
		Addr:         ":" + conf.Port,
	}

	http.HandleFunc("/", authenticate(conf, dispatcher, controller.UpgradeHandler))
	http.HandleFunc("/v1/send", authenticate(conf, dispatcher, v1.Send))
	server.ListenAndServe()
	runtime.Goexit()
}
