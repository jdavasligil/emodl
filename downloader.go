// TODO:
// - Create universal interface for all emotes? Or keep each type separate?
// - Obtain emote image data as well?
// - Provide a stream (iterator) for digesting emote images in batches

package emotedownloader

import (
	"errors"
	"image"
	"io"
	"iter"
	"log"
	"net/http"
	"net/url"
	"slices"
	"sync"

	"golang.org/x/image/webp"

	"github.com/mailru/easyjson"
)

type ImageScale string

const (
	ImageScale1X ImageScale = "1x"
	ImageScale2X            = "2x"
	ImageScale3X            = "3x"
)

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

// TODO:
// - Create with options to supply providers and broadcaster_id
// - Obtain broadcaster custom emotes as well
type EmoteDownloader struct {
	BTTVEmotes
	Err error
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
		// TODO: Reduce allocations here
		ed.BTTVEmotes = slices.Concat(emotes)
	}

	return err
}

// TODO: proper pooling and decoding of webp images asynchronously

// Returns a single-use iterator yielding batches of images for emotes.
// Errors can be checked by calling .Err() on the EmoteDownloader.
func (ed *EmoteDownloader) BTTVImages(imageScale ImageScale) iter.Seq2[string, []image.Image] {
	return func(yield func(string, []image.Image) bool) {
		// TODO: use waitgroup to batch async calls and iterate in chunks
		// We are currently running requests in series which is very slow
		for _, e := range ed.BTTVEmotes {
			id := e.ID
			batch := make([]image.Image, 0, 2)
			r, err := getBTTVEmoteImageData(id, imageScale)
			if err != nil {
				ed.Err = err
				return
			}
			img, err := webp.Decode(r)
			if err != nil {
				ed.Err = err
				return
			}

			r.Close()

			batch = append(batch, img)
			if batch == nil {
				ed.Err = errors.New("Batch can not be nil.")
				return
			}

			if !yield(e.ID, batch) {
				return
			}
		}
		// obtain batch of images in bulk
		// keep track of current emote idx
		// track error in ed
	}
}

// Caller is responsible for closing reader
func getBTTVEmoteImageData(imageID string, imageScale ImageScale) (io.ReadCloser, error) {
	// Example URL to access BTTV cdn for image:
	// https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x.webp
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   "cdn.betterttv.net",
			Path:   "/emote/" + imageID + "/" + string(imageScale) + ".webp",
		},
		Header: http.Header{},
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Error.Message)
	}

	return response.Body, nil
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
