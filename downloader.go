// TODO:
// - Create universal interface for all emotes
// - Asynchronously make requests to each provider
// - Obtain emote image data as well?
// - Provide a stream (iterator) for digesting emotes

package emotedownloader

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x.webp
var (
	bttvAPIVersion = "3"
)

// Used to unmarshal errors from an API response
type jsonError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}
type EmoteDownloader struct {
	Emotes []Emote
}

// BTTV format for now.
type Emote struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId"`
}

func (ed *EmoteDownloader) GetBTTVGlobalEmotes() error {
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
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = json.Unmarshal(body, errorMessage)
		if err != nil {
			return err
		}
		return errors.New(errorMessage.Error.Message)
	}

	// TEMP unmarshal directly into downloader synchronously
	err = json.Unmarshal(body, &ed.Emotes)
	if err != nil {
		return err
	}

	return nil
}
