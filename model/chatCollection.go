package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const chatsFilename = "chats-collection.json"

// Chats is an exported 'Chats' structure
var Chats chats

// chats struct represents a slice of chat objects
type chats struct {
	Collection []Chat `json:"chats"`
}

// Chat object represents a Telegram user, bot, group or a channel.
// We're only interested in its ID for now.
type Chat struct {
	ID int64 `json:"chat_id"`
}

// Recipient implementation of telebot's interface.
// Returns chat_id as string.
func (c Chat) Recipient() string {
	return strconv.FormatInt(c.ID, 10)
}

// Load JSON file into collections
func (db *chats) Load() {
	var err error
	var file *os.File

	if file, err = os.Open(chatsFilename); os.IsNotExist(err) {
		// File does not exist. Create one
		_, err = os.Create(chatsFilename)
		checkFatalError(err)
		return
	}

	// File exists
	fileInfo, err := file.Stat()
	checkFatalError(err)

	if fileInfo.Size() == 0 {
		// Empty file; do not attemp to read it.
		return
	}

	data, err := ioutil.ReadFile(chatsFilename)
	checkFatalError(err)

	err = json.Unmarshal(data, db)
	checkFatalError(err)
}

// Add a user to the collection
func (db *chats) Add(chatID int64) error {
	u := Chat{ID: chatID}

	for _, user := range db.Collection {
		if user == u {
			return errors.New("duplicated ID")
		}
	}

	db.Collection = append(db.Collection, u)
	return nil
}

// Remove a user from the collection
func (db *chats) Remove(chatID int64) {
	u := Chat{ID: chatID}

	log.Printf("User with ID %d wants to be removed...", u.ID)
	for i, user := range db.Collection {
		if user == u {
			copy(db.Collection[i:], db.Collection[i+1:])
			db.Collection[len(db.Collection)-1] = Chat{}
			db.Collection = db.Collection[:len(db.Collection)-1]
		}
	}
}

// Save JSON file.
func (db *chats) Save() {
	data, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
	}

	buf, err := ioutil.ReadFile(chatsFilename)
	if err != nil {
		log.Println(err)
		return
	}

	// No users or groups subscribed since last time
	if bytes.Equal(buf, data) {
		return
	}

	os.Remove(chatsFilename)
	ioutil.WriteFile(chatsFilename, data, 0644)
}
