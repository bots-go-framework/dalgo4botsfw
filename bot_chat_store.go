package dalgo4botsfw

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
)

var _ botsfwdal.BotChatStore = (*botChatStore)(nil)

func newBotChatStore(collection string, getDb DbProvider, newChatData func(botID string) (botsfwmodels.BotChat, error)) botChatStore {
	if collection == "" {
		panic("collection is empty")
	}
	if getDb == nil {
		panic("getDb is nil")
	}
	if newChatData == nil {
		panic("newChatData is nil")
	}
	return botChatStore{
		dalgoStore: dalgoStore{
			getDb:      getDb,
			collection: collection,
		},
		newData: newChatData,
	}
}

type botChatStore struct {
	dalgoStore
	newData func(botID string) (botsfwmodels.BotChat, error)
}

func (store *botChatStore) chatKey(botID, chatID string) *dal.Key {
	return dal.NewKeyWithID(store.collection, fmt.Sprintf("%s:%s", botID, chatID))
}

func (store *botChatStore) GetBotChatEntityByID(c context.Context, botID, chatID string) (botsfwmodels.BotChat, error) {
	key := store.chatKey(botID, chatID)
	data, err := store.newData(botID)
	if err != nil {
		return nil, err
	}
	record := dal.NewRecordWithData(key, data)
	var db dal.Database

	if db, err = store.getDb(c, botID); err != nil {
		return nil, err
	}
	var getter dal.Getter = db
	if tx, ok := dal.GetTransaction(c).(dal.ReadwriteTransaction); ok && tx != nil {
		getter = tx
	}
	if err = getter.Get(c, record); err != nil {
		if dal.IsNotFound(err) {
			err = botsfwdal.NotFoundErr(err)
		}
		return nil, err
	}
	return data, nil
}

func (store *botChatStore) SaveBotChat(c context.Context, botID, chatID string, data botsfwmodels.BotChat) error {
	key := store.chatKey(botID, chatID)
	record := dal.NewRecordWithData(key, data)
	return store.runReadwriteTransaction(c, botID, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, record)
	})
}

func (store *botChatStore) Close(c context.Context) error {
	return nil
}
