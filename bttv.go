package emotedownloader

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mailru/easyjson"
)

var (
	// BTTV EMOTE URL
	// https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x.webp
	bttvAPIVersion = "3"
	bttvHost       = "api.betterttv.net"
)

//easyjson:json
type BTTVEmoteSlice []BTTVEmote

type BTTVEmote struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId"`
}

func getGlobalBTTVEmotes() ([]BTTVEmote, error) {
	sb := strings.Builder{}
	err := apiPathTmpl.Execute(&sb, apiPath{
		Version: bttvAPIVersion,
		Path:    "/cached/emotes/global",
	})
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   bttvHost,
			Path:   sb.String(),
		},
		Header: http.Header{},
	}

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
		err = easyjson.Unmarshal(body, errorMessage)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorMessage.Error.Message)
	}

	bttvBTTVEmotes := BTTVEmoteSlice{}
	err = easyjson.Unmarshal(body, &bttvBTTVEmotes)
	if err != nil {
		return nil, err
	}

	return bttvBTTVEmotes, nil
}
