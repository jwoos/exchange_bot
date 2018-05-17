package main


import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)


func main() {
	log.Print("Beginning application")

	apiToken := os.Getenv("SLACK_EXCHANGE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("API token required")
	}

	verificationToken := os.Getenv("SLACK_EXCHANGE_VERIFICATION_TOKEN")
	if verificationToken == "" {
		log.Fatal("Verification token required")
	}

	client := slack.New(apiToken)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		buffer, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error reading body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body := string(buffer)

		log.Printf("body: %s", body)

		event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{verificationToken}))
		if err != nil {
			log.Printf("error: %v", err)
			log.Print("error parsing event")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch event.Type {
		case slackevents.URLVerification:
			var request slackevents.ChallengeResponse
			err = json.Unmarshal(buffer, &request)
			if err != nil {
				log.Printf("error: %v", err)
				log.Print("error unmarshalling")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(request.Challenge))

		case slackevents.CallbackEvent:
			innerEvent := event.InnerEvent

			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				// TODO check for format
				command := strings.Join(strings.Split(ev.Text, " ")[1:], " ")
				_, _, err := client.PostMessage(ev.Channel, command, slack.PostMessageParameters{})
				if err != nil {
					log.Printf("error: %v", err)
					log.Print("error posting message")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

			case *slackevents.MessageEvent:
				/* FIXME
				 * The library currently has no sub_type so I can't tell
				 * whether the messages are from the bot or user
				 */
				log.Printf("%v", ev)
				_, _, err := client.PostMessage(ev.Channel, ev.Text, slack.PostMessageParameters{})
				if err != nil {
					log.Printf("error: %v", err)
					log.Print("error posting message")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Write(nil)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	const port string = ":8000"
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
