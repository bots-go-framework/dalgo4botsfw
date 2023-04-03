package dalgo4botsfw

import (
	"context"
	"fmt"
	"github.com/strongo/bots-framework/botsfw"
	"github.com/strongo/dalgo/dal"
	"github.com/strongo/dalgo/record"
)

var _ botsfw.BotUserStore = (*botUserStore)(nil)

type BotUserCreator func(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error)

type botUserStore struct {
	dalgoStore
	newBotUserData func() botsfw.BotUser
	createBotUser  BotUserCreator
}

// NewBotUserStore creates new bot user store
func NewBotUserStore(collection string, db DbProvider, newBotUserData func() botsfw.BotUser, createBotUser BotUserCreator) botsfw.BotUserStore {
	if db == nil {
		panic("db is nil")
	}
	if collection == "" {
		panic("collection is empty")
	}
	if newBotUserData == nil {
		panic("newBotUserData is nil")
	}
	return &botUserStore{
		dalgoStore: dalgoStore{
			db:         db,
			collection: collection,
		},
		newBotUserData: newBotUserData,
		createBotUser:  createBotUser,
	}
}

type botUserWithInt64ID struct {
	record.WithID[int64]
	Data botsfw.BotUser
}

// GetBotUserByID returns bot user data
func (store botUserStore) GetBotUserByID(c context.Context, botUserID any) (botsfw.BotUser, error) {
	key := store.botUserKey(botUserID)
	botUserData := store.newBotUserData()
	botUser := botUserWithInt64ID{
		Data: botUserData,
		WithID: record.WithID[int64]{
			ID:     botUserID.(int64),
			Record: dal.NewRecordWithData(key, botUserData),
		},
	}
	db, err := store.db(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}
	return botUser.Data, db.Get(c, botUser.Record)
}

// SaveBotUser saves bot user data
func (store botUserStore) SaveBotUser(c context.Context, botUserID any, botUserData botsfw.BotUser) error {
	key := store.botUserKey(botUserID)
	record := dal.NewRecordWithData(key, botUserData)
	db, err := store.db(c)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Set(c, record)
	})
}

func (store botUserStore) CreateBotUser(c context.Context, botID string, apiUser botsfw.WebhookActor) (botsfw.BotUser, error) {
	return store.createBotUser(c, botID, apiUser)
}

func (store botUserStore) botUserKey(botUserID any) *dal.Key {
	return dal.NewKeyWithID(store.collection, botUserID.(int64))
}
