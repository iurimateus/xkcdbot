package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

const helpMsg string = `Options:
/current        - Send the current comic
/comic (num)    - Send comic #num. If not specified, sends a random comic
/subscribe      - receive new comics
/unsubscribe    - stop receiving new comics`

const usageMsg string = "Usage: /comic {number}.\nExample: `comic 1000`"

// execEvery executes a function repeatedly at a specific interval
// see: reddit.com/r/golang/comments/8auh9j/is_there_a_way_to_supervise_and_restart/dx2w5k1/
func execEvery(interval time.Duration, f func()) {
	tick := time.NewTimer(interval)
	for range tick.C {
		f()
		tick.Reset(interval)
	}
}

func checkBotError(m *tb.Message, err error) bool {
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "checkBotError"))

		if m.Chat != nil {
			bot.Send(m.Chat, err.Error())
		}
		return true
	}

	return false
}
