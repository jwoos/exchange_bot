package main


import (
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)


var serverLogger = initializeLogger("server")


type Server struct {
	router *mux.Router
	client *slack.Client
	token *Token
	users map[string]*User
	context *slack.AuthTestResponse
}


func newServer() *Server {
	// get token
	token := newToken()

	client := slack.New(token.api)

	context, err := client.AuthTest()
	if err != nil {
		serverLogger.Fatalf("Failed to authenticate: %v", err)
	}

	server := &Server{
		router: mux.NewRouter().StrictSlash(true),
		token: token,
		client: client,
		users: make(map[string]*User),
		context: context,
	}

	return server
}
