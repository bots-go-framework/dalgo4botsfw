package dalgo4botsfw

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwdal"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

var _ botsfwdal.BotUserStore = (*botUserStore)(nil)

type BotUserCreator func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfwmodels.BotUser, error)

type botUserStore struct {
	dalgoStore
	newBotUserData func(botID string) (botsfwmodels.BotUser, error)
	createBotUser  BotUserCreator
}

// newBotUserStore creates new bot user store
func newBotUserStore(collection string, getDb DbProvider, newBotUserData func(botID string) (botsfwmodels.BotUser, error), createBotUser BotUserCreator) botUserStore {
	if getDb == nil {
		panic("getDb is nil")
	}
	if collection == "" {
		panic("collection is empty")
	}
	if newBotUserData == nil {
		panic("newBotUserData is nil")
	}
	return botUserStore{
		dalgoStore: dalgoStore{
			getDb:      getDb,
			collection: collection,
		},
		newBotUserData: newBotUserData,
		createBotUser:  createBotUser,
	}
}

type botUserWithStrID struct {
	record.WithID[string]
	Data botsfwmodels.BotUser
}

// GetBotUserByID returns bot user data
func (store botUserStore) GetBotUserByID(c context.Context, botID, botUserID string) (botsfwmodels.BotUser, error) {
	key := store.botUserKey(botUserID)
	botUserData, err := store.newBotUserData(botID)
	if err != nil {
		return nil, err
	}
	botUser := botUserWithStrID{
		Data: botUserData,
		WithID: record.WithID[string]{
			ID:     botUserID,
			Record: dal.NewRecordWithData(key, botUserData),
		},
	}
	db, err := store.getDb(c, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get getDb: %w", err)
	}

	var getter dal.Getter = db
	if tx, ok := dal.GetTransaction(c).(dal.ReadwriteTransaction); ok && tx != nil {
		getter = tx
	}

	if err = getter.Get(c, botUser.Record); err != nil {
		if dal.IsNotFound(err) {
			err = botsfwdal.NotFoundErr(err)
		}
		return nil, err
	}
	return botUser.Data, nil
}

// SaveBotUser saves bot user data
func (store botUserStore) SaveBotUser(c context.Context, botID, botUserID string, botUserData botsfwmodels.BotUser) error {
	key := store.botUserKey(botUserID)
	botUserRecord := dal.NewRecordWithData(key, botUserData)
	return store.runReadwriteTransaction(c, botID, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, botUserRecord)
	})
}

func (store botUserStore) CreateBotUser(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfwmodels.BotUser, error) {
	return store.createBotUser(c, botID, apiUser)
}

func (store botUserStore) botUserKey(botUserID any) *dal.Key {
	switch id := botUserID.(type) {
	case string:
		return dal.NewKeyWithID(store.collection, id)
	case int:
		return dal.NewKeyWithID(store.collection, id)
	case int64:
		return dal.NewKeyWithID(store.collection, id)
	default:
		panic(fmt.Sprintf("unsupported botUserID type: %T: %v", botUserID, botUserID))
	}
}
