package dalgo4botsfw

import (
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
)

var _ botsfwdal.DataAccess = (*dataAccess)(nil)

type dataAccess struct {
	getDb        DbProvider
	recordsMaker botsfwmodels.BotRecordsMaker
	botChatStore
	botUserStore
	appUserStore
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
	getDb DbProvider,
	recordsMaker botsfwmodels.BotRecordsMaker,
) botsfwdal.DataAccess {
	if getDb == nil {
		panic("getDb == nil")
	}
	if recordsMaker == nil {
		panic("recordsMaker == nil")
	}
	return &dataAccess{
		getDb:        getDb,
		botChatStore: newBotChatStore("botChat", getDb, recordsMaker.MakeBotChatDto),
		botUserStore: newBotUserStore("botUser", getDb, recordsMaker.MakeBotUserDto, nil),
		appUserStore: newAppUserStore("appUser", getDb),
	}
}
