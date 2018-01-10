package dispatcher

import (
	"net/http"

	"github.com/upsub/dispatcher/src/config"
)

// Listen starts the dispatcher
func Listen() {
	config := config.Create()
	dispatcher := createDispatcher()
	go dispatcher.serve()

	http.HandleFunc("/", authenticate(config, dispatcher, upgradeHandler))
	http.ListenAndServe(":"+config.Port, nil)
}
