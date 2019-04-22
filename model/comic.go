package model

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	// Use current Unix time as seed
	rand.Seed(time.Now().UTC().UnixNano())
}

// Comic JSON structure.
type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link,omitempty"`
	Year       string `json:"year"`
	News       string `json:"news,omitempty"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript,omitempty"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

// FetchLastComic fetches the latest in-memory XKCD comic
func FetchLastComic() (*Comic, error) {
	return FetchComic(Comics.Last)
}

// FetchRandomComic returns a pointer to a random comic
func FetchRandomComic() (*Comic, error) {
	upperBound := Comics.Last

	num := rand.Intn(upperBound)
	return FetchComic(num)
}

// FetchComic requires a number and returns a pointer to a Comic struct.
func FetchComic(num int) (*Comic, error) {
	comic, ok := Comics.Collection[num]

	if !ok {
		numStr := strconv.Itoa(num)
		urlSuffix := "/info.0.json"

		req, err := Client.newRequest(http.MethodGet, numStr+urlSuffix)
		if err != nil {
			return nil, err
		}

		_, err = Client.do(req, &comic, false)
		if err != nil {
			return nil, err
		}

		Comics.Add(&comic, false)
	}

	return &comic, nil
}

// FetchRecentComic fetches the current xkcd comic.
func FetchRecentComic() (*Comic, error) {
	var comic Comic

	req, err := Client.newRequest(http.MethodGet, "/info.0.json")
	if err != nil {
		return nil, err
	}

	res, err := Client.do(req, &comic, true)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPermanentRedirect {
		return nil, fmt.Errorf(res.Status)
	}

	return &comic, nil
}
