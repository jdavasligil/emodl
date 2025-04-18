// TODO:
// - Create universal interface for all emotes
// - Obtain emote image data as well?
// - Provide a stream (iterator) for digesting emotes

package emotedownloader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"
)

// Example URL to access BTTV cdn for image:
// https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x.webp
var (
	bttvAPIVersion = "3"

	EmoteProviders = []string{
		"BTTV",
		"7TV",
		"FFZ",
	}
)

// Used to unmarshal errors from an API response
type jsonError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// BTTV format for now.
type Emote struct {
	Provider  string
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId"`
}

type EmoteDownloader struct {
	Emotes []Emote
}

// Download and collect the emote data from each provider.
func (ed *EmoteDownloader) Download() error {
	var err error

	emotesChan := make(chan []Emote, 4)
	errorChan := make(chan error, 4)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		bttvEmotes, err := getBTTVGlobalEmotes()
		emotesChan <- bttvEmotes
		errorChan <- err
	}()

	// Block until all requests finish.
	wg.Wait()

	// Channels must be closed to iterate over their contents.
	close(emotesChan)
	close(errorChan)

	for e := range errorChan {
		err = errors.Join(e)
	}

	for emotes := range emotesChan {
		ed.Emotes = slices.Concat(emotes)
	}

	return err
}

// Returns a function which can be repeatedly called to obtain a batch of emote images.
func (ed *EmoteDownloader) ImageIterator(buf *bytes.Buffer) func() {
	return func() {
		// obtain batch of images in bulk
		// keep track of current emote idx
		// track error in ed
		// write image data to buffer
		// how to properly deserialize images based on type?
		// associate emote id to image?
	}
}

func getBTTVGlobalEmotes() ([]Emote, error) {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.betterttv.net",
			Path:   "/" + bttvAPIVersion + "/cached/emotes/global",
		},
		Header: http.Header{},
	}

	//req.Header.Set("Client-ID", b.ClientID)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = json.Unmarshal(body, errorMessage)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Error.Message)
	}

	bttvEmotes := []Emote{}
	err = json.Unmarshal(body, &bttvEmotes)
	if err != nil {
		return nil, err
	}

	for i, _ := range bttvEmotes {
		bttvEmotes[i].Provider = "BTTV"
	}

	return bttvEmotes, nil
}
