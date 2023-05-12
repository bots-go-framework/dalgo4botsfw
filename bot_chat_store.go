package dalgo4botsfw

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
)

var _ botsfwdal.BotChatStore = (*botChatStore)(nil)

func newBotChatStore(collection, platform string, getDb DbProvider, newChatData func(botID string) (botsfwmodels.ChatData, error)) botChatStore {
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
		platform: platform,
		newData:  newChatData,
	}
}

type botChatStore struct {
	dalgoStore
	platform string
	newData  func(botID string) (botsfwmodels.ChatData, error)
}

func (store *botChatStore) botChatRecordKey(k botsfwmodels.ChatKey) *dal.Key {
	return dal.NewKeyWithID(store.collection, fmt.Sprintf("%s:%s:%s", store.platform, k.BotID, k.ChatID))
}

func (store *botChatStore) GetBotChatData(c context.Context, chatKey botsfwmodels.ChatKey) (chatData botsfwmodels.ChatData, err error) {
	recordKey := store.botChatRecordKey(chatKey)
	chatData, err = store.newData(chatKey.BotID)
	chatData.Base().ChatKey = chatKey
	if err != nil {
		return
	}
	var db dal.Database

	if db, err = store.getDb(c, chatKey.BotID); err != nil {
		return nil, err
	}

	var getter dal.Getter = db

	if tx := dal.GetTransaction(c); tx != nil {
		var isRW bool
		if getter, isRW = tx.(dal.ReadTransaction); !isRW {
			getter = nil
		}
	}
	record := dal.NewRecordWithData(recordKey, chatData)
	if err = getter.Get(c, record); err != nil {
		if dal.IsNotFound(err) {
			err = botsfwdal.NotFoundErr(err)
		}
	}
	return
}

func (store *botChatStore) SaveBotChatData(c context.Context, chatKey botsfwmodels.ChatKey, data botsfwmodels.ChatData) error {
	key := store.botChatRecordKey(chatKey)
	record := dal.NewRecordWithData(key, data)
	return store.runReadwriteTransaction(c, chatKey.BotID, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, record)
	})
}

func (store *botChatStore) Close(c context.Context) error {
	return nil
}
