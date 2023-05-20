package dalgo4botsfw

import (
	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
)

var _ botsfwdal.DataAccess = (*dataAccess)(nil)

type dataAccess struct {
	getDb DbProvider
	botChatStore
	botUserStore
	appUserStore
	//recordsMaker botsfwmodels.BotRecordsMaker
}

func (da dataAccess) RunInTransaction(c context.Context, botID string, f func(c context.Context) error) error {
	db, err := da.getDb(c, botID)
	if err != nil {
		return err
	}
	return db.RunReadwriteTransaction(c, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return f(ctx)
	})
}

func NewDataAccess(
	platform string,
	getDb DbProvider,
	recordsMaker botsfwmodels.BotRecordsMaker,
) botsfwdal.DataAccess {
	if getDb == nil {
		panic("getDb == nil")
	}
	getDbWrapper := func(c context.Context, botID string) (dal.Database, error) {
		db, err := getDb(c, botID)
		if err != nil {
			return db, err
		}
		if db == nil {
			return nil, errors.New("db provider returned no db")
		}
		return db, nil
	}
	if recordsMaker == nil {
		panic("recordsMaker == nil")
	}
	return &dataAccess{
		getDb:        getDb,
		botChatStore: newBotChatStore("botChat", platform, getDbWrapper, recordsMaker.MakeBotChatDto),
		botUserStore: newBotUserStore("botUser", platform, getDbWrapper, recordsMaker.MakeBotUserDto, nil),
		appUserStore: newAppUserStore("appUser", getDbWrapper),
	}
}
