package main


type User struct {
	id string
	money int
	portfolio *Portfolio
}


func newUser(id string) *User {
	user := &User{
		id: id,
		money: CONFIG_MONEY_BEGIN,
		portfolio: newPortfolio(),
	}

	return user
}

func getOrCreateUser(users map[string]*User, id string) *User {
	val, ok := users[id]
	if !ok {
		user := newUser(id)
		users[id] = user
		val = user
	}

	return val
}
