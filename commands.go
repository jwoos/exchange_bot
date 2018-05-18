package main


import (
	"strings"
)


var commandMap = map[string]func(*Server, *User) (string, error){
	"help": helpCommand,
}

var helpMap = map[string]string{
	"help" : "View this dialogue",
	"get <symbol>": "Get information for <symbol> where it is a stock or cryptocurrency symbol",
	"buy <symbol> <amount>": "Buy <amount> of <symbol>",
	"sell <symbol> <amount>": "Sell <amount> of <symbol>",
	"balance": "View your available balance",
	"portfolio": "View your portfolio",
	"leaderboard": "View the leaderboard",
	"greet": "Return a greeting",
	"meme": "Post a meme",
}


func helpCommand(s *Server, u *User) (string, error) {
	builder := strings.Builder{}

	var err error
	for k, v := range helpMap {
		_, err = builder.WriteString(strings.Join([]string{k, v}, ": "))
		if err != nil {
			return "", err
		}

		_, err = builder.WriteRune('\n')
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}
