package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"xkcdbot/model"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	model.Chats.Load()
	model.Comics.Load()
}

func main() {
	go bot.Start()
	defer bot.Stop()

	// Save data to disk before closing...
	defer model.Chats.Save()
	defer model.Comics.Save()

	// Ensure a update will not be sent on the first time polling
	if model.Comics.Last == 0 {
		// empty or missing collections file.
		comic, err := model.FetchRecentComic()
		if err != nil {
			log.Panicln(err)
		}

		model.Comics.Add(comic, true)
	}

	// Polls most recent comic every five minutes
	go execEvery(5*time.Minute, poller)

	// Write collection to disk every hour
	go execEvery(time.Hour, model.Comics.Save)

	// Stop on Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("\n\nShutting down...")
}
