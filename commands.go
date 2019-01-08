package main


import (
	"fmt"
	"strings"

	"github.com/jwoos/slack_exchange/assets"
)


var commandsLogger = initializeLogger("commands")


var commandMap = map[string]func(*Server, *User, []string) (string, error){
	"help": helpCommand,
	"price": priceCommand,
	//"quote": quoteCommand,
	"balance": balanceCommand,
}

var helpMap = map[string]string{
	"help" : "View this dialog",
	"price (c[rypto|s[tock]) <symbol> ...": "Get price for <symbol>",
	"quote (c[rypto|s[tock]) <symbol> ...": "Get quote for <symbol>",
	"buy (c[rypto|s[tock]) <symbol> <amount>": "Buy <amount> of <symbol>",
	"sell (c[rypto]|s[tock]) <symbol> <amount>": "Sell <amount> of <symbol>",
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
	builder.WriteString("Look below for a valid command \n===============\n")
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

func priceCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	symbols := make([]string, len(cmd[2:]))
	for i, sym := range cmd[2:] {
		symbols[i] = strings.ToUpper(sym)
	}

	var err error
	switch cmd[1] {
	case "s":
		fallthrough
	case "stock":
		iex := assets.IEXMarketBatch{}
		err = iex.Fetch(assets.IEXRequest{
			Information: []string{"price"},
			Symbols: symbols,
		})
		if err != nil {
			commandsLogger.Errorf("error fetching stock price: %v", err)
		}

		for sym, to := range iex.Batch {
			builder.WriteString(fmt.Sprintf("%s: %f\n", sym, *to.Price))
		}

	case "c":
		fallthrough
	case "crypto":
		cc := assets.CCMulti{}
		err = cc.Fetch(assets.CCRequest{
			FromSymbols: symbols,
			ToSymbols: []string{"USD"},
		})
		if err != nil {
			commandsLogger.Errorf("error fetching crypto price: %v", err)
		}

		for sym, to := range cc.Batch {
			usd, _ := to["USD"]
			builder.WriteString(fmt.Sprintf("%s: %f\n", sym, usd))
		}

	default:
		builder.WriteString("Invalid option, please give one of s[tock] or c[rypto]")
	}

	return builder.String(), err
}

/*
 *func quoteCommand(s *Server, u *User, cmd []string) (string, error) {
 *    builder := strings.Builder{}
 *
 *    symbols := make([]string, len(cmd[2:]))
 *    for i, sym := range cmd[2:] {
 *        symbols[i] = strings.ToUpper(sym)
 *    }
 *
 *    var err error
 *    switch cmd[1] {
 *    case "s":
 *        fallthrough
 *    case "stock":
 *        iex := assets.IEXMarketBatch{}
 *        err = iex.Fetch(assets.IEXRequest{
 *            Information: []string{"price"},
 *            Symbols: symbols,
 *        })
 *        if err != nil {
 *            commandsLogger.Errorf("error fetching stock price: %v", err)
 *        }
 *
 *        for sym, to := range iex.Batch {
 *            builder.WriteString(fmt.Sprintf("%s: %f\n", sym, *to.Price))
 *        }
 *
 *    case "c":
 *        fallthrough
 *    case "crypto":
 *        cc := assets.CCMulti{}
 *        err = cc.Fetch(assets.CCRequest{
 *            FromSymbols: symbols,
 *            ToSymbols: []string{"USD"},
 *        })
 *        if err != nil {
 *            commandsLogger.Errorf("error fetching crypto price: %v", err)
 *        }
 *
 *        for sym, to := range cc.Batch {
 *            usd, _ := to["USD"]
 *            builder.WriteString(fmt.Sprintf("%s: %f\n", sym, usd))
 *        }
 *
 *    default:
 *        builder.WriteString("Invalid option, please give one of s[tock] or c[rypto]")
 *    }
 *
 *    return builder.String(), err
 *}
 */

func balanceCommand(s *Server, u *User, cmd []string) (string, error) {
	balance := fmt.Sprintf("%d", u.money)
	return balance, nil
}

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
