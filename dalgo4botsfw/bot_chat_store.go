package dalgo4botsfw

import (
	"context"
	"github.com/strongo/bots-framework/botsfw"
)

var _ botsfw.BotChatStore = (*botChatStore)(nil)

type botChatStore struct {
	dalgoStore
}

func (b botChatStore) GetBotChatEntityByID(c context.Context, botID, botChatID string) (botsfw.BotChat, error) {
	//TODO implement me
	panic("implement me")
}

func (b botChatStore) SaveBotChat(c context.Context, botID, botChatID string, chatEntity botsfw.BotChat) error {
	//TODO implement me
	panic("implement me")
}

func (b botChatStore) NewBotChatEntity(c context.Context, botID string, botChat botsfw.WebhookChat, appUserID, botUserID string, isAccessGranted bool) botsfw.BotChat {
	//TODO implement me
	panic("implement me")
}

func (b botChatStore) Close(c context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewBotChatStore(collection string, db DbProvider) botsfw.BotChatStore {
	if collection == "" {
		panic("collection is empty")
	}
	if db == nil {
		panic("db is nil")
	}
	return botChatStore{
		dalgoStore: dalgoStore{
			db:         db,
			collection: collection,
		},
	}
}
