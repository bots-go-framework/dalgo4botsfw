package dalgo4botsfw

import (
	"github.com/strongo/bots-framework/botsfw"
	"github.com/strongo/dalgo/record"
)

type botChat struct {
	record.WithID[string]
	Data botsfw.BotChat
}
