package emodl

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"unsafe"

	"github.com/mailru/easyjson"
)

var (
	bttvAPIVersion     = "3"
	bttvHost           = "api.betterttv.net"
	bttvCDNPathTmpl, _ = template.New("bttvCDN").Parse("https://cdn.betterttv.net/emote/{{ .ID }}/1x.webp")
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
	UserID    string `json:"userID,intern"`
}

func (e BTTVEmote) URL() string {
	var url strings.Builder
	err := bttvCDNPathTmpl.Execute(&url, e)
	if err != nil {
		panic(err)
	}
	return url.String()
}

func (e BTTVEmote) Image() Image {
	return Image{
		URL:    e.URL(),
		Width:  28,
		Height: 28,
		ID:     e.ID,
	}
}

func (e BTTVEmote) AsEmote() Emote {
	return Emote{
		ID: e.ID,
		Name: e.Name,
		Images: []Image{e.Image()},
		Locations: []string{},
	}
}

func (e BTTVEmote) Size() uintptr {
	return unsafe.Sizeof(e)
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
