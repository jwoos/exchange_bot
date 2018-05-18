package main


import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)


func urlVerificationEvent(s *Server, w http.ResponseWriter, buffer []byte) {
	var request slackevents.ChallengeResponse
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		log.Printf("error: %v", err)
		log.Print("error unmarshalling")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(request.Challenge))
}

func callbackEvent(s *Server, w http.ResponseWriter, event slackevents.EventsAPIInnerEvent) {
	switch ev := event.Data.(type) {
	case *slackevents.AppMentionEvent:
		user := getOrCreateUser(s.users, ev.User)

		// TODO check for format
		command := strings.Split(ev.Text, " ")[1:]

		fn, okay := commandMap[command[0]]
		if !okay {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Invalid Command"))
			return
		}

		response, err := fn(s, user)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		_, _, err = s.client.PostMessage(
			ev.Channel,
			response,
			slack.PostMessageParameters{},
		)
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error posting message")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case *slackevents.MessageEvent:
		user := getOrCreateUser(s.users, ev.User)

		// TODO check for format
		command := strings.Split(ev.Text, " ")

		_, _, err := s.client.PostMessage(
			ev.Channel,
			fmt.Sprintf("%s %s %d", strings.Join(command, " "), user.id, user.money),
			slack.PostMessageParameters{},
		)
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error posting message")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (s *Server) handleEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error reading body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := string(buffer)

		log.Printf("body: %s", body)

		event, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionVerifyToken(&slackevents.TokenComparator{s.token.verification}),
		)
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error parsing event")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch event.Type {
		case slackevents.URLVerification:
			urlVerificationEvent(s, w, buffer)

		case slackevents.CallbackEvent:
			callbackEvent(s, w, event.InnerEvent)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
