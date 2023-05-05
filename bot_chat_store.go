package dalgo4botsfw

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
)

var _ botsfw.BotChatStore = (*botChatStore)(nil)

func NewBotChatStore(collection string, db DbProvider, newData func() botsfw.BotChat) botsfw.BotChatStore {
	if collection == "" {
		panic("collection is empty")
	}
	if db == nil {
		panic("db is nil")
	}
	if newData == nil {
		panic("newData is nil")
	}
	return botChatStore{
		dalgoStore: dalgoStore{
			db:         db,
			collection: collection,
		},
		newData: newData,
	}
}

type botChatStore struct {
	dalgoStore
	newData func() botsfw.BotChat
}

func (store botChatStore) chatKey(botID, chatID string) *dal.Key {
	return dal.NewKeyWithID(store.collection, fmt.Sprintf("%s:%s", botID, chatID))
}

func (store botChatStore) GetBotChatEntityByID(c context.Context, botID, chatID string) (botsfw.BotChat, error) {
	key := store.chatKey(botID, chatID)
	data := store.newData()
	record := dal.NewRecordWithData(key, data)
	db, err := store.db(c)
	if err != nil {
		return nil, err
	}
	return data, db.Get(c, record)
}

func (store botChatStore) SaveBotChat(c context.Context, botID, chatID string, data botsfw.BotChat) error {
	db, err := store.db(c)
	if err != nil {
		return err
	}
	key := store.chatKey(botID, chatID)
	record := dal.NewRecordWithData(key, data)
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, record)
	})
}

func (store botChatStore) NewBotChatEntity(c context.Context, botID string, botChat botsfw.WebhookChat, appUserID, botUserID string, isAccessGranted bool) botsfw.BotChat {
	panic("implement me") //TODO implement me
}

func (store botChatStore) Close(c context.Context) error {
	panic("implement me") //TODO implement me
}
