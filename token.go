package main


import (
	"os"
)


var tokenLogger = initializeLogger("token")


type Token struct {
	api string
	verification string
}


func newToken() *Token {
	apiToken := os.Getenv("SLACK_EXCHANGE_API_TOKEN")
	if apiToken == "" {
		tokenLogger.Fatal("API token required")
	}

	verificationToken := os.Getenv("SLACK_EXCHANGE_VERIFICATION_TOKEN")
	if verificationToken == "" {
		tokenLogger.Fatal("Verification token required")
	}

	token := &Token{
		api: apiToken,
		verification: verificationToken,
	}

	return token
}
