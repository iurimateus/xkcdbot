package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"xkcdbot/model"

	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

func init() {
	settings := tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	}

	var err error
	bot, err = tb.NewBot(settings)
	if err != nil {
		log.Fatalln(err)
	}

	// Bot functionality
	bot.Handle("/start", Start)
	bot.Handle("/help", help)
	bot.Handle("/stop", stop)

	// Comic
	bot.Handle("/current", current)
	bot.Handle("/comic", comicByNumber)

	// Receive updates?
	bot.Handle("/subscribe", subHandler)
	bot.Handle("/unsubscribe", unsubHandler)
}

// Start handler
func Start(m *tb.Message) {
	if !m.Private() {
		return
	}

	help(m)
}

func help(m *tb.Message) {
	bot.Send(m.Chat, helpMsg)
}

// Add it to file and force a update
func subscribe(chat *tb.Chat) error {
	err := model.Chats.Add(chat.ID)
	if err == nil {
		model.Chats.Save()
	}

	return err
}

// Remove from file and force a update
func unsubscribe(chat *tb.Chat) {
	model.Chats.Remove(chat.ID)
	model.Chats.Save()
}

func subHandler(m *tb.Message) {
	var msg string

	err := subscribe(m.Chat)
	if err != nil {
		if err.Error() == "duplicated ID" {
			msg = "You're already subscribed!"
		} else {
			log.Println(err)
			msg = err.Error()
		}

		bot.Send(m.Chat, msg)
		return
	}

	bot.Send(m.Chat, "Subscribed")
}

func unsubHandler(m *tb.Message) {
	unsubscribe(m.Chat)

	bot.Send(m.Chat, "Subscription cancelled")
}

func stop(m *tb.Message) {
	unsubHandler(m)
}

func randomComic(m *tb.Message) {
	comic, err := model.FetchRandomComic()
	if checkBotError(m, err) {
		return
	}

	sendComic(m, comic)
}

func current(m *tb.Message) {
	comic, err := model.FetchLastComic()
	if checkBotError(m, err) {
		return
	}

	sendComic(m, comic)
}

func comicByNumber(m *tb.Message) {
	if len(m.Payload) == 0 {
		// Empty message.
		// Send a random comic
		randomComic(m)
		return
	}

	id, err := strconv.Atoi(m.Payload)
	if err != nil {
		if strings.Contains(err.Error(), "invalid syntax") {
			// invalid command
			bot.Send(m.Chat, usageMsg, tb.ModeMarkdown)
			return
		}

		// Unknown error. Log and send it to the user
		checkBotError(m, err)
		return
	}

	if id > model.Comics.Last || id < 0 {
		// id out of range
		bot.Send(m.Chat, "Non-existent number... sending current comic")
		id = model.Comics.Last
	}

	comic, err := model.FetchComic(id)
	if checkBotError(m, err) {
		return
	}

	sendComic(m, comic)
}

func sendComic(m *tb.Message, comic *model.Comic) {
	chat := model.Chat{ID: m.Chat.ID}

	if !m.Private() {
		// Enable retries for public chats (group/channel)
		sendComicToChat(chat, comic, true, 0)
		return
	}

	sendComicToChat(chat, comic, false, 0)
}

func sendComicToChat(chat model.Chat, comic *model.Comic, retry bool, retries uint) {
	// Display as: [1500] Upside-Down Map
	msg := fmt.Sprintf("[%d] %s\n\n", comic.Num, comic.Title)

	// Exponential backoff
	const maxRetries = 5
	if retries <= maxRetries {
		m, err := bot.Send(chat, &tb.Photo{File: tb.FromURL(comic.Img), Caption: msg + comic.Alt})
		if !checkBotError(m, err) || !retry {
			return
		}

		log.Println("Error: ", err)

		msg := err.Error() + " Trying again..."
		bot.Send(chat, msg)

		// 2 ^ retries * 150 milliseconds
		exp := (1 << retries) * 150
		totalWait := exp + rand.Intn(1000) // random time in ms up to a second
		sleepTime := time.Duration(totalWait) * time.Millisecond

		log.Printf("exp: %v,totalWait: %v", exp, totalWait)

		log.Printf("Try: %d, sleeping for %s\n", retries, sleepTime.String())
		time.Sleep(sleepTime)

		go sendComicToChat(chat, comic, true, retries+1)
	}
}

// newComic sends latest comic to all subscribed chats
func newComic() {
	key := model.Comics.Last
	last := model.Comics.Collection[key]

	log.Println("New Comic!", last.Num, last.Title)

	coll := model.Chats.Collection
	for _, chat := range coll {
		go sendComicToChat(chat, &last, true, 0)
	}
}
