package dalgo4botsfw

import "github.com/strongo/dalgo/dal"

type dalgoStore struct {
	collection string
	db         dal.Database
}
