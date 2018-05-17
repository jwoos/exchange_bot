package main


import (
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)


type Server struct {
	router *mux.Router
	client *slack.Client
	token *Token
}


func newServer() *Server {
	// get token
	token := newToken()

	server := &Server{
		router: mux.NewRouter().StrictSlash(true),
		token: token,
		client: slack.New(token.api),
	}

	return server
}
