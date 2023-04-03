package botsfw2dalgo

import (
	"github.com/strongo/bots-framework/botsfw"
	"github.com/strongo/dalgo/dal"
	"testing"
)

func TestNewBotUserStore(t *testing.T) {
	type args struct {
		db             dal.Database
		collection     string
		newBotUserData func() botsfw.BotUser
	}
	tests := []struct {
		name        string
		args        args
		shouldPanic bool
	}{
		{name: "empty", shouldPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewBotUserStore() did not panic")
					}
				}()
			}
			if got := NewBotUserStore(tt.args.db, tt.args.collection, tt.args.newBotUserData); got == nil {
				t.Error("NewBotUserStore() returned nil")
			}
		})
	}
}
