package dalgo4botsfw

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var _ botsfw.BotAppUserStore = (*appUserStore)(nil)

type appUserStore struct {
	dalgoStore
}

func (a appUserStore) GetAppUserByID(c context.Context, appUserID string, appUser botsfw.BotAppUser) error {
	panic("implement me") //TODO implement me
}

func (a appUserStore) CreateAppUser(c context.Context, botID string, actor botsfw.WebhookActor) (appUserID string, appUserEntity botsfw.BotAppUser, err error) {
	panic("implement me") //TODO implement me
}

func NewAppUserStore(collection string, db DbProvider) botsfw.BotAppUserStore {
	if collection == "" {
		panic("collection is empty")
	}
	if db == nil {
		panic("db is nil")
	}
	return appUserStore{
		dalgoStore: dalgoStore{
			db:         db,
			collection: collection,
		},
	}
}
