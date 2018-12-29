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
	"github.com/upsub/dispatcher/server/message"
)

var (
	conf *config.Config
	dis  *dispatcher.Dispatcher
)

// Listen starts the http server and creates a new dispatcher
func Listen() {
	conf = config.Create()
	store := auth.NewStore(conf)
	dis = dispatcher.Create(conf, store)
	go dis.Serve()

	server := &http.Server{
		ReadTimeout:  conf.ConnectionTimeout,
		WriteTimeout: conf.ConnectionTimeout,
		Addr:         ":" + strconv.Itoa(int(conf.Port)),
	}

	http.HandleFunc("/", authenticate(conf, dis, store, controller.UpgradeHandler))
	http.HandleFunc("/v1/send", authenticate(conf, dis, store, v1.Send))
	server.ListenAndServe()
	runtime.Goexit()
}

// Send broadcasts a message to all listening clients
func Send(msg *message.Message) {
	dis.ProcessMessage(msg, dispatcher.CreateConnection(
		"dispatcher",
		"",
		"dispatcher",
		"1.0",
		nil,
		conf,
		dis,
		map[string]bool{
			"wildcard": true,
		},
	))
}
