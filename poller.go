package main

import (
	"log"
	"xkcdbot/model"
)

// poller polls the API and sends a update to the channel if the response contains a newer comic
func poller() {
	comic, err := model.FetchRecentComic()
	if err != nil {
		log.Println(err)
		return
	}

	if comic.Num > model.Comics.Last {
		// Polled is newer!
		model.Comics.Add(comic, true)
		model.Comics.Save() // Force file update

		newComic()
	}
}
