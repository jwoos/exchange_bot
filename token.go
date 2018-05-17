package main


import (
	"log"
	"os"
)


type Token struct {
	api string
	verification string
}


func newToken() *Token {
	apiToken := os.Getenv("SLACK_EXCHANGE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("API token required")
	}

	verificationToken := os.Getenv("SLACK_EXCHANGE_VERIFICATION_TOKEN")
	if verificationToken == "" {
		log.Fatal("Verification token required")
	}

	token := &Token{
		api: apiToken,
		verification: verificationToken,
	}

	return token
}
