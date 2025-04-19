package emotedownloader

// DOCUMENTATION
// https://github.com/SevenTV/API/blob/main/internal/api/rest/v3/docs/swagger.json
// https://cdn.7tv.app/emote/01F6MQ33FG000FFJ97ZB8MWV52/3x.avif
// https://7tv.io/v3/emote-sets/global

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/mailru/easyjson"
)

var (
	sevenTVAPIVersion = "v3"
	sevenTVHost       = "7tv.io"

	sevenTVPathTmpl, _ = template.New("api").Parse("/{{ .Version }}/{{ .Path }}/{{ .Option }}")
)

type sevenTVPath struct {
	Version string
	Path    string
	Option  string
}

//easyjson:json
type SevenTVEmoteCollection struct {
	Name   string `json:"name"`
	Emotes []struct {
		Data SevenTVEmote `json:"data"`
	} `json:"emotes"`
}

type SevenTVEmote struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Animated bool   `json:"animated"`
	Host     struct {
		Url   string `json:"url,intern"`
		Files []struct {
			Name       string `json:"name"`
			StaticName string `json:"static_name"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			FrameCount int    `json:"frame_count"`
			Size       uint32 `json:"size"`
			Format     string `json:"format"`
		} `json:"files"`
	} `json:"host"`
}

// Performs GET request for the emote collection
func get7TVEmoteCollection(collection string) (*SevenTVEmoteCollection, error) {
	sb := strings.Builder{}
	err := sevenTVPathTmpl.Execute(&sb, sevenTVPath{
		Version: sevenTVAPIVersion,
		Path:    "emote-sets",
		Option:  collection,
	})
	if err != nil {
		return nil, err
	}
	//fmt.Println(sb.String())

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "https",
			Host:   sevenTVHost,
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

	c := &SevenTVEmoteCollection{}
	err = easyjson.UnmarshalFromReader(response.Body, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
