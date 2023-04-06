package dalgo4botsfw

import (
	"context"
	"github.com/dal-go/dalgo/dal"
)

type DbProvider func(c context.Context) (dal.Database, error)

type dalgoStore struct {
	collection string
	db         DbProvider
}
