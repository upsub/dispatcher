package dispatcher

import (
	"net/http"
)

// Listen starts the dispatcher
func Listen() {
	config := createConfig()
	dispatcher := createDispatcher()
	go dispatcher.serve()

	http.HandleFunc("/", authenticate(config, dispatcher))
	http.ListenAndServe(":"+config.port, nil)
}
