package main

import (
	"os"
)

var tokenLogger = initializeLogger("token")

type Token struct {
	api          string
	signing      string
	verification string
}

func newToken() *Token {
	apiToken := os.Getenv("SLACK_EXCHANGE_API_TOKEN")
	if apiToken == "" {
		tokenLogger.Fatal("API token required")
	}

	// new verification token
	signingToken := os.Getenv("SLACK_EXCHANGE_SIGNING_TOKEN")
	/*
	 *if signingToken == "" {
	 *    tokenLogger.Fatal("Signing token required")
	 *}
	 */

	// old verification token
	verificationToken := os.Getenv("SLACK_EXCHANGE_VERIFICATION_TOKEN")
	tokenLogger.Warning("Use signing token instead if possible")
	if verificationToken == "" {
		tokenLogger.Fatal("Verification token required")
	}

	token := &Token{
		api:          apiToken,
		signing:      signingToken,
		verification: verificationToken,
	}

	return token
}
