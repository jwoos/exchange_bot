package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/jwoos/slack_exchange/assets"

	"github.com/nlopes/slack"
)

var commandsLogger = initializeLogger("commands")

var commandMap = map[string]func(*Server, *User, []string) (string, error){
	"help":  helpCommand,
	"price": priceCommand,
	//"quote": quoteCommand,
	"buy": buyCommand,
	//"sell": sellCommand,
	"balance":     balanceCommand,
	"portfolio":   portfolioCommand,
	"leaderboard": leaderboardCommand,
	//"greet": greatCommand,
	//"meme": memeCommand
}

var helpMap = map[string]string{
	"help":                                      "View this dialog",
	"price (c[rypto|s[tock]) <symbol> ...":      "Get price for <symbol>",
	"quote (c[rypto|s[tock]) <symbol> ...":      "Get quote for <symbol>",
	"buy (c[rypto|s[tock]) <symbol> <amount>":   "Buy <amount> of <symbol>",
	"sell (c[rypto]|s[tock]) <symbol> <amount>": "Sell <amount> of <symbol>",
	"balance":     "View your available balance",
	"portfolio":   "View your portfolio",
	"leaderboard": "View the leaderboard",
	"greet":       "Return a greeting",
	"meme":        "Post a meme",
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
			Symbols:     symbols,
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
			ToSymbols:   []string{"USD"},
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

	return builder.String(), nil
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

 func buyCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	symbol := strings.ToUpper(cmd[2])
	count, err := (strconv.ParseFloat(cmd[3], 64))
	if err != nil {
		return "Couldn't parse amount as a number", nil
	}

	var price float64
	var stock bool

	switch cmd[1] {
	case "s":
		fallthrough
	case "stock":
		stock = true

		count = math.Floor(count)

		iex := assets.IEXMarketBatch{}
		err = iex.Fetch(assets.IEXRequest{
			Information: []string{"price"},
			Symbols:     []string{symbol},
		})
		if err != nil {
			commandsLogger.Errorf("error fetching stock price: %v", err)
			return "Error fetching prices - try again later", err
		}

		to, okay := iex.Batch[symbol]
		if !okay {
			return fmt.Sprintf("Symbol %s was not found", symbol), nil
		}

		price = *to.Price

	case "c":
		fallthrough
	case "crypto":
		stock = false

		cc := assets.CCMulti{}
		err = cc.Fetch(assets.CCRequest{
			FromSymbols: []string{symbol},
			ToSymbols:   []string{"USD"},
		})
		if err != nil {
			commandsLogger.Errorf("error fetching crypto price: %v", err)
			return "Error fetching prices - try again later", err
		}

		to, okay := cc.Batch[symbol]
		if !okay {
			return fmt.Sprintf("Symbol %s was not found", symbol), nil
		}

		price = to["USD"]

	default:
		builder.WriteString("Invalid option, please give one of s[tock] or c[rypto]")
		return builder.String(), nil
	}

	total := price * count
	if total > u.balance {
		return fmt.Sprintf("You request to buy %.2f but you have %.2f", u.balance, total), nil
	}

	u.balance -= total

	if (stock) {
		u.portfolio.appendStock(newAsset(symbol, price, count))
	} else {
		u.portfolio.appendCrypto(newAsset(symbol, price, count))
	}
	builder.WriteString(fmt.Sprintf("Bought %.2f of %s @ %.2f\n", count, symbol, price))

	return builder.String(), nil
 }

/*
 * func sellCommand(s *Server, u *User, cmd []string) (string, error) {
 *
 * }
 */

func balanceCommand(s *Server, u *User, cmd []string) (string, error) {
	balance := fmt.Sprintf("%.2f", u.balance)
	return balance, nil
}

func leaderboardCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	stockCache := make(map[string]float64)
	cryptoCache := make(map[string]float64)

	// get all symbols
	for _, user := range s.users {
		for k, _ := range user.portfolio.stock {
			stockCache[k] = 0
		}

		for k, _ := range user.portfolio.cryptocurrency {
			cryptoCache[k] = 0
		}
	}

	var err error
	var index uint

	stockKeys := make([]string, len(stockCache))
	index = 0
	for k, _ := range stockCache {
		stockKeys[index] = k
		index++
	}

	iex := assets.IEXMarketBatch{}
	err = iex.Fetch(assets.IEXRequest{
		Symbols:     stockKeys,
		Information: []string{"price"},
	})
	if err != nil {
		commandsLogger.Errorf("error fetching stock price: %v", err)
		return "Error fetching prices - try again later", err
	}

	cryptoKeys := make([]string, len(cryptoCache))
	index = 0
	for k, _ := range cryptoCache {
		cryptoKeys[index] = k
		index++
	}

	cc := assets.CCMulti{}
	err = cc.Fetch(assets.CCRequest{
		FromSymbols: cryptoKeys,
		ToSymbols:   []string{"USD"},
	})
	if err != nil {
		commandsLogger.Errorf("error fetching crypto price: %v", err)
		return "Error fetching prices - try again later", err
	}

	users := make([]struct{
		user *slack.User
		total float64
	}, len(s.users))

	index = 0
	for _, user := range s.users {
		total := user.balance

		for k, _ := range user.portfolio.stock {
			to := iex.Batch[k]
			total += *to.Price
		}

		for k, _ := range user.portfolio.cryptocurrency {
			to := cc.Batch[k]
			total += to["USD"]
		}

		users[index] = struct {
			user *slack.User
			total float64
		}{
			user: user.slackUser,
			total: total,
		}
	}

	sort.Slice(users, func(i int, j int) bool {
		return users[i].total > users[j].total
	})

	builder.WriteString("*Leaderboard*\n")

	index = 0
	for _, user := range users {
		builder.WriteString(fmt.Sprintf(
			"%d - %s (%s): %.2f\n",
			index,
			user.user.Name,
			user.user.RealName,
			user.total,
		))
		index++
	}

	return builder.String(), nil
}

func portfolioCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	portfolio := u.portfolio

	builder.WriteString("*Stock*\n")
	for symbol, assets := range portfolio.stock {
		builder.WriteString(fmt.Sprintf("_%s_\n", symbol))
		for _, asset := range assets {
			builder.WriteString(fmt.Sprintf("%.0f @ %f\n", asset.count, asset.price))
		}
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	builder.WriteString("*Cryptocurrency*\n")
	for symbol, assets := range portfolio.cryptocurrency {
		builder.WriteString(fmt.Sprintf("_%s_\n", symbol))
		for _, asset := range assets {
			builder.WriteString(fmt.Sprintf("%.0f @ %f\n", asset.count, asset.price))
		}
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	return builder.String(), nil
}
