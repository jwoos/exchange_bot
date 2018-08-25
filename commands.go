package main


import (
	"fmt"
	"strings"
)


var commandMap = map[string]func(*Server, *User, []string) (string, error){
	"help": helpCommand,
	"sget": getCommand,
	"cget": getCommand,
	"balance": balanceCommand,
}

// TODO
var helpMap = map[string]string{
	"help" : "View this dialog",
	"sget <symbol>": "Get information for <symbol> stock",
	"cget <symbol>": "Get information for <symbol> crypto",
	"sbuy <symbol> <amount>": "Buy <amount> of <symbol>",
	"cbuy <symbol> <amount>": "Buy <amount> of <symbol>",
	"ssell <symbol> <amount>": "Sell <amount> of <symbol>",
	"csell <symbol> <amount>": "Sell <amount> of <symbol>",
	"balance": "View your available balance",
	"portfolio": "View your portfolio",
	"leaderboard": "View the leaderboard",
	"greet": "Return a greeting",
	"meme": "Post a meme",
}


func errorCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	builder.WriteString(
		fmt.Sprintf("Invalid command: %s\n", strings.Join(cmd, " ")),
	)
	builder.WriteString("Look below for a valid command \n================\n")
	helpDialog, _ := helpCommand(s, u, cmd)
	builder.WriteString(helpDialog)

	return builder.String(), nil
}

func helpCommand(s *Server, u *User, cmd []string) (string, error) {
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

// TODO
func getCommand(s *Server, u *User, cmd []string) (string, error) {
	//builder := strings.Builder{}

	var err error
	switch cmd[1] {
	case "stock":
	case "crypto":
	default:
	}

	return "", err
}

func balanceCommand(s *Server, u *User, cmd []string) (string, error) {
	balance := fmt.Sprintf("%d", u.money)
	return balance, nil
}

// TODO
/*
 *func portfolioCommand(s *Server, u *User, cmd []string) (string, error) {
 *    builder := strings.Builder{}
 *
 *    portfolio := u.portfolio
 *
 *    stock := portfolio.stock
 *    builder.WriteString("Stock\n")
 *    for symbol, assets := range stock {
 *        for asset := range assets {
 *
 *        }
 *    }
 *    cryptocurrency := portfolio.cryptocurrency
 *
 *    return builder.String(), nil
 *}
 */
