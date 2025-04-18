// TODO:
// - Create universal interface for all emotes
// - Obtain emote image data as well?
// - Provide a stream (iterator) for digesting emotes

package emotedownloader

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"sync"

	"github.com/mailru/easyjson"
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
//
//easyjson:json
type jsonError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// BTTV format for now.
//
//easyjson:json
type BTTVEmote struct {
	Provider  string
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType,intern"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId,intern"`
}

//easyjson:json
type BTTVEmotes []BTTVEmote

type EmoteDownloader struct {
	BTTVEmotes
}

// Download and collect the emote data from each provider.
func (ed *EmoteDownloader) Download() error {
	var err error

	emotesChan := make(chan BTTVEmotes, 4)
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
		ed.BTTVEmotes = slices.Concat(emotes)
	}

	return err
}

// Returns a function which can be repeatedly called to obtain a batch of emote images.
func (ed *EmoteDownloader) ImageIterator() func() {
	return func() {
		// obtain batch of images in bulk
		// Security consideration: caution with decoding potentially large images: https://pkg.go.dev/image
		// keep track of current emote idx
		// track error in ed
		// write image data to buffer
		// how to properly deserialize images based on type?
		// associate emote id to image?
	}
}

func getBTTVGlobalEmotes() (BTTVEmotes, error) {
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

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Error.Message)
	}

	bttvEmotes := BTTVEmotes{}
	err = easyjson.UnmarshalFromReader(response.Body, &bttvEmotes)
	if err != nil {
		return nil, err
	}

	for i, _ := range bttvEmotes {
		bttvEmotes[i].Provider = "BTTV"
	}

	return bttvEmotes, nil
}
