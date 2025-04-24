package emodl

// DOCUMENTATION
// https://betterttv.com/developers/api

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"text/template"
	"unsafe"

	"github.com/mailru/easyjson"
)

var (
	bttvAPIVersion      = "3"
	bttvHost            = "api.betterttv.net"
	bttvCDNPathTmpl, _  = template.New("bttvCDN").Parse("https://cdn.betterttv.net/emote/{{ .ID }}/1x.webp")
	bttvUserPathTmpl, _ = template.New("bttvCDN").Parse("/{{ .Version }}/cached/users/{{ .Platform }}/{{ .PlatformID }}")
)

type bttvUserPath struct {
	Version    string
	Platform   string
	PlatformID string
}

// Either BTTVTVID or Platform/PlatformID are needed to get user emote sets.
type BTTVOptions struct {
	// Platform linked to BTTV (twitch, youtube)
	Platform string

	// ID associated with Platform (not username)
	PlatformID string
}

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
	ID   string `json:"id"`
	Name string `json:"code"`
	// The image type is a lie. We always can get webp.
	//ImageType string `json:"imageType,intern"`
	Animated bool `json:"animated"`
	//UserID   string `json:"userID,intern"`
}

func (e BTTVEmote) URL() string {
	var url strings.Builder
	err := bttvCDNPathTmpl.Execute(&url, e)
	if err != nil {
		panic(err)
	}
	return url.String()
}

// We just have to assume all images are 28x28. No way to query other than
// downloading the image directly. To be safe, just get this data from the image.
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
		ID:        e.ID,
		Name:      e.Name,
		Images:    []Image{e.Image()},
		Locations: []string{},
	}
}

func (e BTTVEmote) Size() uintptr {
	return unsafe.Sizeof(e)
}

//easyjson:json
type BTTVUser struct {
	ID            string      `json:"id"`
	ChannelEmotes []BTTVEmote `json:"channelEmotes"`
	SharedEmotes  []BTTVEmote `json:"sharedEmotes"`
}

func getBTTVUser(platform string, platformID string) (BTTVUser, error) {
	var u BTTVUser
	sb := strings.Builder{}
	err := bttvUserPathTmpl.Execute(&sb, bttvUserPath{
		Version:    bttvAPIVersion,
		Platform:   platform,
		PlatformID: platformID,
	})
	if err != nil {
		panic(err)
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
		return u, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return u, err
			}
			return u, errors.New(response.Status + "\n" + string(body))
		}
		return u, errors.New(errorMessage.Error.Message)
	}

	err = easyjson.UnmarshalFromReader(response.Body, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func getBTTVUserEmotes(platform string, platformID string) (BTTVEmoteSlice, error) {
	bttvEmotes := BTTVEmoteSlice{}
	u, err := getBTTVUser(platform, platformID)
	if err != nil {
		return bttvEmotes, err
	}
	return slices.Concat(u.SharedEmotes, u.ChannelEmotes), nil
}

func getBTTVGlobalEmotes() (BTTVEmoteSlice, error) {
	var bttvEmotes BTTVEmoteSlice
	sb := strings.Builder{}
	err := apiPathTmpl.Execute(&sb, apiPath{
		Version: bttvAPIVersion,
		Path:    "/cached/emotes/global",
	})
	if err != nil {
		return bttvEmotes, err
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
		return bttvEmotes, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return bttvEmotes, err
			}
			return bttvEmotes, errors.New(response.Status + "\n" + string(body))
		}
		return bttvEmotes, errors.New(errorMessage.Error.Message)
	}

	err = easyjson.UnmarshalFromReader(response.Body, &bttvEmotes)
	if err != nil {
		return bttvEmotes, err
	}

	return bttvEmotes, nil
}
