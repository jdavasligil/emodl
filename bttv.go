package emotedownloader

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unsafe"

	"github.com/mailru/easyjson"
)

var (
	bttvAPIVersion = "3"
	bttvHost       = "api.betterttv.net"
	// https://cdn.betterttv.net/emote/54fa8f1401e468494b85b537/1x.webp
	bttvCDN = "cdn.betterttv.net"
)

//easyjson:json
type BTTVEmoteSlice []BTTVEmote

func (es BTTVEmoteSlice) Size() uintptr {
	if es == nil {
		return 0
	}
	var size uintptr

	size += unsafe.Sizeof(es)

	for _, e := range es {
		size += e.Size()
	}

	return size
}

type BTTVEmote struct {
	ID        string `json:"id"`
	Name      string `json:"code"`
	ImageType string `json:"imageType,intern"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId,intern"`
}

func (e *BTTVEmote) Size() uintptr {
	if e == nil {
		return 0
	}
	return unsafe.Sizeof(*e)
}

func getGlobalBTTVEmotes() (BTTVEmoteSlice, error) {
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

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(response.Status + "\n" + string(body))
		}
		return nil, errors.New(errorMessage.Error.Message)
	}

	bttvBTTVEmotes := BTTVEmoteSlice{}
	err = easyjson.UnmarshalFromReader(response.Body, &bttvBTTVEmotes)
	if err != nil {
		return nil, err
	}

	return bttvBTTVEmotes, nil
}
