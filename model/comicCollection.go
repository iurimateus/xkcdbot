package model

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const comicsFilename = "comics-collection.json"

// Comics exports an empty 'comics' struct
var Comics = &comics{
	Collection: make(map[int]Comic),
}

type comics struct {
	Collection map[int]Comic `json:"comics"`
	Last       int           `json:"last"`
}

// Load JSON file into collections
func (db *comics) Load() {
	var err error
	var file *os.File

	if file, err = os.Open(comicsFilename); os.IsNotExist(err) {
		// File does not exist. Create one
		_, err = os.Create(comicsFilename)
		checkFatalError(err)
		return
	}

	// File exists
	fInfo, err := file.Stat()
	checkFatalError(err)

	if fInfo.Size() == 0 {
		// Empty file; do not attempt to read it.
		return
	}

	data, err := ioutil.ReadFile(comicsFilename)
	checkFatalError(err)

	err = json.Unmarshal(data, db)
	checkFatalError(err)
}

// Save JSON file.
func (db *comics) Save() {
	data, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := ioutil.ReadFile(comicsFilename)
	if err != nil {
		log.Println(err)
		return
	}

	if bytes.Equal(buf, data) {
		// File matches in-memory data
		return
	}

	os.Remove(comicsFilename)
	ioutil.WriteFile(comicsFilename, data, 0644)
}

// Add requires a pointer to 'comic' and a boolean. Adds the comic to a collection.
func (db *comics) Add(comic *Comic, latest bool) {
	key := comic.Num
	Comics.Collection[key] = *comic

	if latest {
		Comics.Last = key
	}
}
