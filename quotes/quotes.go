package quotes

import (
	"errors"
	"math/rand"
	"time"
)

//Quotes is a map that stores username mapped to an array of quotes
type Quotes map[string][]string

//New creates a new map for storing user quotes
func New() Quotes {
	return make(map[string][]string)
}

//Add adds a quote attributed to a user
func (q Quotes) Add(user string, quote string) {
	if _, ok := q[user]; !ok {
		q[user] = []string{quote}
	} else {
		for existingUser, quotes := range q {
			if user == existingUser {
				quotes = append(quotes, quote)
			}
		}
	}
}

//Get returns a random quote from a user
func (q Quotes) Get(user string) (string, error) {
	var quote string
	if len(q) == 0 {
		return quote, errors.New("No quotes have been added for any users")
	}
	for storedUser, quotes := range q {
		if user == storedUser {
			if len(quotes) == 1 {
				quote = quotes[0]
				return quote, nil
			}
			rand.Seed(time.Now().UnixNano())
			quote = quotes[rand.Intn(len(quotes)-1)]
			return quote, nil
		}
	}
	return quote, errors.New("User has no quotes")
}
