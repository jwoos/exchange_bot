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
}


func newServer() *Server {
	// get token
	token := newToken()

	server := &Server{
		router: mux.NewRouter().StrictSlash(true),
		token: token,
		client: slack.New(token.api),
		users: make(map[string]*User),
	}

	_, err := server.client.AuthTest()
	if err != nil {
		serverLogger.Fatalf("Failed to authenticate: %v", err)
	}

	return server
}
