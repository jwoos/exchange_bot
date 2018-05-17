package main


import (
	"log"
	"net/http"
)


func main() {
	// set up server
	server := newServer()

	// register routes
	server.routes()

	const port string = ":8000"
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, server.router))
}
