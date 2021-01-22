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
	"sell": sellCommand,
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

	if len(cmd) < 3 {
		return "Invalid number of arguments - look at help", nil
	}

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
			builder.WriteString(fmt.Sprintf("%s: %.2f\n", sym, *to.Price))
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
			builder.WriteString(fmt.Sprintf("%s: %.2f\n", sym, usd))
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
 *            builder.WriteString(fmt.Sprintf("%s: %.2f\n", sym, *to.Price))
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
 *            builder.WriteString(fmt.Sprintf("%s: %.2f\n", sym, usd))
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

	if len(cmd) != 4 {
		return "Invalid number of arguments - look at help", nil
	}

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
		return fmt.Sprintf("You request to buy %.0f but you have %.2f", u.balance, total), nil
	}

	u.balance -= total

	if (stock) {
		u.portfolio.addStock(symbol, count)
	} else {
		u.portfolio.addCrypto(symbol, count)
	}
	builder.WriteString(fmt.Sprintf("Bought %.0f of %s @ %.2f for %.2f\n", count, symbol, price, total))

	return builder.String(), nil
 }

 func sellCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	if len(cmd) != 4 {
		return "Invalid number of arguments - look at help", nil
	}

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

		_, ok := u.portfolio.stock[symbol]
		if !ok {
			return fmt.Sprintf("Symbol %s was not found in your stocks portfolio", symbol), nil
		}

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

		_, ok := u.portfolio.cryptocurrency[symbol]
		if !ok {
			return fmt.Sprintf("Symbol %s was not found in your cryptocurrency portfolio", symbol), nil
		}

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

	if (stock) {
		_, err := u.portfolio.removeStock(symbol, count)
		if err != nil {
			return err.Error(), nil
		}
		builder.WriteString(fmt.Sprintf("Sold %.0f of %s @ %.2f\n", count, symbol, price))
	} else {
		_, err := u.portfolio.removeCrypto(symbol, count)
		if err != nil {
			return err.Error(), nil
		}
		builder.WriteString(fmt.Sprintf("Sold %.2f of %s @ %.2f\n", count, symbol, price))
	}

	total := price * count
	u.balance += total

	return builder.String(), nil
 }

func balanceCommand(s *Server, u *User, cmd []string) (string, error) {
	balance := fmt.Sprintf("%.2f", u.balance)
	return balance, nil
}

func leaderboardCommand(s *Server, u *User, cmd []string) (string, error) {
	builder := strings.Builder{}

	stockSet := make(map[string]struct{})
	cryptoSet := make(map[string]struct{})

	// get all symbols
	for _, user := range s.users {
		for k, _ := range user.portfolio.stock {
			stockSet[k] = struct{}{}
		}

		for k, _ := range user.portfolio.cryptocurrency {
			cryptoSet[k] = struct{}{}
		}
	}

	var err error
	var index uint
	var iex assets.IEXMarketBatch
	var cc assets.CCMulti

	if len(stockSet) != 0 {
		stockKeys := make([]string, len(stockSet))
		index = 0
		for k, _ := range stockSet {
			stockKeys[index] = k
			index++
		}

		err = iex.Fetch(assets.IEXRequest{
			Symbols:     stockKeys,
			Information: []string{"price"},
		})
		if err != nil {
			commandsLogger.Errorf("error fetching stock price: %v", err)
			return "Error fetching prices - try again later", err
		}
	}

	if len(cryptoSet) != 0 {
		cryptoKeys := make([]string, len(cryptoSet))
		index = 0
		for k, _ := range cryptoSet {
			cryptoKeys[index] = k
			index++
		}

		err = cc.Fetch(assets.CCRequest{
			FromSymbols: cryptoKeys,
			ToSymbols:   []string{"USD"},
		})
		if err != nil {
			commandsLogger.Errorf("error fetching crypto price: %v", err)
			return "Error fetching prices - try again later", err
		}
	}

	users := make([]struct{
		user *slack.User
		total float64
	}, len(s.users))

	index = 0
	for _, user := range s.users {
		total := user.balance

		for k, v := range user.portfolio.stock {
			to := iex.Batch[k]
			total += *to.Price * v.count
		}

		for k, v := range user.portfolio.cryptocurrency {
			to := cc.Batch[k]
			total += to["USD"] * v.count
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

	builder.WriteString(fmt.Sprintf("*Balance*\n%.2f\n\n", u.balance))

	builder.WriteString("*Stock*\n")
	for symbol, asset := range u.portfolio.stock {
		builder.WriteString(fmt.Sprintf("_%s_ - %.0f\n", symbol, asset.count))
	}
	builder.WriteString("\n")

	builder.WriteString("*Cryptocurrency*\n")
	for symbol, asset := range u.portfolio.cryptocurrency {
		builder.WriteString(fmt.Sprintf("_%s_ - %.2f\n", symbol, asset.count))
	}
	builder.WriteString("\n")

	return builder.String(), nil
}
