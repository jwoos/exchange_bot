package main


import (
	"github.com/nlopes/slack"
)


type User struct {
	id string
	slackUser *slack.User
	balance float64
	portfolio *Portfolio
}


func newUser(slackUser *slack.User) *User {
	user := &User{
		id: slackUser.ID,
		slackUser: slackUser,
		balance: CONFIG_MONEY_BEGIN,
		portfolio: newPortfolio(),
	}

	return user
}
