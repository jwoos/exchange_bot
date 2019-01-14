package main

import (
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)

var serverLogger = initializeLogger("server")

type Server struct {
	router  *mux.Router
	client  *slack.Client
	token   *Token
	users   map[string]*User
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
		router:  mux.NewRouter().StrictSlash(true),
		token:   token,
		client:  client,
		users:   make(map[string]*User),
		context: context,
	}

	return server
}

func (s *Server) getOrCreateUser(id string) (*User, error) {
	user, ok := s.users[id]
	if !ok {
		slackUser, err := s.client.GetUserInfo(id)
		if err != nil {
			serverLogger.Errorf("Failed to fetch information for user %s: %v", id, err)
			return nil, err
		}

		user = newUser(slackUser)
		s.users[id] = user

		serverLogger.Infof("User created: %v", user)
	}

	return user, nil
}

func (s *Server) sendMessage(channel string, message string) error {
	_, _, err := s.client.PostMessage(
		channel,
		message,
		slack.PostMessageParameters{
			Markdown: true,
		},
	)

	if err != nil {
		errMessage := fmt.Errorf("Failed sending message to Slack: %v", err)
		serverLogger.Error(errMessage)
		return errMessage
	}

	return nil
}
