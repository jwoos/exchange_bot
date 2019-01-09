package main

import (
	"fmt"
	"net/http"
)

var mainLogger = initializeLogger("main")

func main() {
	// set up server
	server := newServer()

	// register routes
	server.routes()

	mainLogger.Infof("Listening on %s", PORT)
	err := http.ListenAndServe(PORT, server.router)
	mainLogger.Critical(fmt.Sprintf("%v", err))
}
