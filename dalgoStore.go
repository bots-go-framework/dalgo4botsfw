package dalgo4botsfw

import (
	"context"
	"github.com/dal-go/dalgo/dal"
)

type DbProvider func(c context.Context, botID string) (dal.Database, error)

type dalgoStore struct {
	collection string
	getDb      DbProvider
}

func (s dalgoStore) runReadwriteTransaction(c context.Context, botID string, f func(c context.Context, tx dal.ReadwriteTransaction) error) error {
	db, err := s.getDb(c, botID)
	if err != nil {
		return err
	}
	return db.RunReadwriteTransaction(c, f)
}
