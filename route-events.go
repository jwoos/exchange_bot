package main


import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)


var routeEventsLogger = initializeLogger("route-events")


func urlVerificationEvent(s *Server, w http.ResponseWriter, buffer []byte) {
	var request slackevents.ChallengeResponse
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		routeEventsLogger.Warningf("Error unmarshalling: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(request.Challenge))
}

func callbackEvent(s *Server, w http.ResponseWriter, event slackevents.EventsAPIInnerEvent) {
	var slackUser string
	var slackChannel string
	var slackText string

	switch event.Type {
	case slackevents.AppMention:
		ev := event.Data.(*slackevents.AppMentionEvent)

		slackUser = ev.User
		slackChannel = ev.Channel
		slackText = ev.Text
	case slackevents.Message:
		ev := event.Data.(*slackevents.MessageEvent)

		slackUser = ev.User
		slackChannel = ev.Channel
		slackText = ev.Text
	default:
		// Error here
	}

	user := getOrCreateUser(s.users, slackUser)

	// TODO check for format
	// Normalize to lowercase
	command := strings.Split(strings.ToLower(slackText), " ")[1:]

	fn, okay := commandMap[command[0]]
	if !okay {
		response, _ := errorCommand(s, user, command)
		_, _, err := s.client.PostMessage(
			slackChannel,
			response,
			slack.PostMessageParameters{},
		)

		if err != nil {
			routeEventsLogger.Errorf("Failed sending message: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	response, err := fn(s, user, command)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, _, err = s.client.PostMessage(
		slackChannel,
		response,
		slack.PostMessageParameters{},
	)
	if err != nil {
		routeEventsLogger.Errorf("Failed sending message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (s *Server) handleEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer, err := ioutil.ReadAll(r.Body)
		if err != nil {
			routeEventsLogger.Errorf("error reading body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := string(buffer)

		event, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionVerifyToken(&slackevents.TokenComparator{s.token.verification}),
		)
		if err != nil {
			routeEventsLogger.Errorf("error parsing event: %v", err)
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
