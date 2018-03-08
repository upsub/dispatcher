package server

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/upsub/dispatcher/server/auth"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/controller"
	"github.com/upsub/dispatcher/server/controller/v1"
	"github.com/upsub/dispatcher/server/dispatcher"
)

// Listen starts the http server and creates a new dispatcher
func Listen() {
	conf := config.Create()
	store := auth.NewStore(conf)
	dispatcher := dispatcher.Create(conf, store)
	go dispatcher.Serve()

	server := &http.Server{
		ReadTimeout:  conf.ConnectionTimeout,
		WriteTimeout: conf.ConnectionTimeout,
		Addr:         ":" + strconv.Itoa(int(conf.Port)),
	}

	http.HandleFunc("/", authenticate(conf, dispatcher, store, controller.UpgradeHandler))
	http.HandleFunc("/v1/send", authenticate(conf, dispatcher, store, v1.Send))
	server.ListenAndServe()
	runtime.Goexit()
}
