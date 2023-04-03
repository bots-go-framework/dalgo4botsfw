package dalgo4botsfw

import (
	"github.com/strongo/dalgo/dal"
	"testing"
)

func TestNewBotChatStore(t *testing.T) {
	type args struct {
		collection string
		db         dal.Database
	}
	tests := []struct {
		name        string
		args        args
		shouldPanic bool
	}{
		{name: "empty", args: args{collection: "", db: nil}, shouldPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewBotChatStore() did not panic")
					}
				}()
			}
			if got := NewBotChatStore(tt.args.collection, tt.args.db); got == nil {
				t.Error("NewBotChatStore() returned nil")
			}
		})
	}
}
