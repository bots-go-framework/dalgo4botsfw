package dalgo4botsfw

import (
	"context"
	"github.com/strongo/dalgo/dal"
)

type DbProvider func(c context.Context) (dal.Database, error)

type dalgoStore struct {
	collection string
	db         DbProvider
}
