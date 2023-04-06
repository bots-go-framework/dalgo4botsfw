package dalgo4botsfw

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/record"
)

type botChat struct {
	record.WithID[string]
	Data botsfw.BotChat
}
