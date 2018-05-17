package main


import (
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"os"

	//"github.com/nlopes/slack"
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

	//api := slack.New(apiToken)

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

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	const port string = ":8000"
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
