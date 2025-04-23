package emodl

// DOCUMENTATION
// https://github.com/SevenTV/API/blob/main/internal/api/rest/v3/docs/swagger.json
// https://www.openapiviewer.com/view/2025-04-23T16-45-41-003Z--566eac48a5217ad2c1fb50a6cd8ccf6e87671e225590eb7f1a9f1c0d8c8586e6--swagger-json#?route=overview
//
// URL EXAMPLES
// https://cdn.7tv.app/emote/01F6MQ33FG000FFJ97ZB8MWV52/3x.avif
// https://7tv.io/v3/emote-sets/global

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
	sevenTVAPIVersion = "v3"
	sevenTVHost       = "7tv.io"

	sevenTVPathTmpl, _ = template.New("sevenTVPath").Parse("/{{ .Version }}/{{ .Path }}/{{ .Option }}")
)

type sevenTVPath struct {
	Version string
	Path    string
	Option  string
}

// Either SevenTVID or Platform/PlatformID are needed to get user emote sets.
type SevenTVOptions struct {
	// Platform linked to 7TV (Twitch, YouTube, Discord)
	Platform string

	// ID associated with Platform (not username)
	PlatformID string

	// ID associated with 7TV directly
	SevenTVID string
}

//easyjson:json
type SevenTVUser struct {
	ID        string `json:"id"`
	EmoteSets []struct {
		ID string `json:"id"`
	} `json:"emote_sets"`
}

//easyjson:json
type SevenTVPlatformUser struct {
	User SevenTVUser `json:"user"`
}

//easyjson:json
type SevenTVEmoteSet struct {
	Name   string `json:"name"`
	Emotes []struct {
		Data SevenTVEmote `json:"data"`
	} `json:"emotes"`
}

func (c SevenTVEmoteSet) Size() uintptr {
	var size uintptr

	size += unsafe.Sizeof(c)
	for _, e := range c.Emotes {
		size += e.Data.Size()
	}

	return size
}

type SevenTVEmote struct {
	ID       string `json:"id"`
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

// GetImage returns the data of an image of a given size and format.
//
// Will always attempt to provide a correct image url nearest to intention.
// Failure to obtain desired scale or format leads to attempted fallbacks.
//
// @param scale  The intended scale of the image ("1x", "2x", etc.)
// @param format The image filetype requested (webp, png, gif, etc.)
func (e *SevenTVEmote) GetImage(scale string, format string) (Image, error) {
	var img Image
	if e == nil {
		return img, errors.New("7TV Emote is nil")
	} else if e.Host.Files == nil {
		return img, errors.New("7TV Emote host files are nil")
	} else if len(e.Host.Files) == 0 {
		return img, errors.New("7TV Emote has no host files")
	} else if e.Host.Url == "" {
		return img, errors.New("7TV Emote has no host url")
	}
	var url strings.Builder
	var imgID strings.Builder

	imgID.WriteString(e.ID)
	imgID.WriteByte('+')

	url.WriteString("https:")
	url.WriteString(e.Host.Url)
	url.WriteByte('/')

	// Find scale / format match
	format = strings.ToUpper(format)

	for _, f := range e.Host.Files {
		if format == f.Format && f.Name[:2] == scale {
			url.WriteString(f.Name)
			imgID.WriteString(f.Name)
			img.ID = imgID.String()
			img.Height = f.Height
			img.Width = f.Width
			img.URL = url.String()
			return img, nil
		}
	}

	// If no match was found, check fallback formats
	for _, fallback := range imageFallbacks {
		for _, f := range e.Host.Files {
			if fallback == f.Format && f.Name[:2] == scale {
				url.WriteString(f.Name)
				imgID.WriteString(f.Name)
				img.ID = imgID.String()
				img.Height = f.Height
				img.Width = f.Width
				img.URL = url.String()
				return img, nil
			}
		}
	}

	// If no match was still found, simply return the first thing we find
	url.WriteString(e.Host.Files[0].Name)
	imgID.WriteString(e.Host.Files[0].Name)
	img.ID = imgID.String()
	img.Height = e.Host.Files[0].Height
	img.Width = e.Host.Files[0].Width
	img.URL = url.String()

	return img, nil
}

func (e *SevenTVEmote) AsEmote() (Emote, error) {
	img, err := e.GetImage("1x", "WEBP")
	if err != nil {
		return Emote{}, err
	}
	return Emote{
		ID: e.ID,
		Name: e.Name,
		Images: []Image{img},
		Locations: []string{},
	}, err
}

func (e *SevenTVEmote) Size() uintptr {
	if e == nil {
		return 0
	}
	var size uintptr

	size += unsafe.Sizeof(*e)

	for _, file := range e.Host.Files {
		size += unsafe.Sizeof(file)
	}

	return size
}

func get7TVEmoteSet(setid string) (SevenTVEmoteSet, error) {
	sb := strings.Builder{}
	set := SevenTVEmoteSet{}
	err := sevenTVPathTmpl.Execute(&sb, sevenTVPath{
		Version: sevenTVAPIVersion,
		Path:    "emote-sets",
		Option:  setid,
	})
	if err != nil {
		return set, err
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
		return set, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := &jsonError{}
		err = easyjson.UnmarshalFromReader(response.Body, errorMessage)
		if err != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return set, err
			}
			return set, errors.New(response.Status + "\n" + string(body))
		}
		return set, errors.New(errorMessage.Error.Message)
	}

	err = easyjson.UnmarshalFromReader(response.Body, &set)
	if err != nil {
		return set, err
	}

	return set, nil
}

func get7TVUser(uid string) (SevenTVUser, error) {
	sb := strings.Builder{}
	u := SevenTVUser{}
	err := sevenTVPathTmpl.Execute(&sb, sevenTVPath{
		Version: sevenTVAPIVersion,
		Path:    "users",
		Option:  uid,
	})
	if err != nil {
		return u, err
	}

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

func get7TVPlatformUser(platform string, pid string) (SevenTVUser, error) {
	sb := strings.Builder{}
	u := SevenTVUser{}
	err := sevenTVPathTmpl.Execute(&sb, sevenTVPath{
		Version: sevenTVAPIVersion,
		Path:    "users/" + platform,
		Option:  pid,
	})
	if err != nil {
		return u, err
	}

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

	pu := SevenTVPlatformUser{}
	err = easyjson.UnmarshalFromReader(response.Body, &pu)
	if err != nil {
		return u, err
	}
	return pu.User, nil

}

func get7TVUserEmoteSetIDs(opt *SevenTVOptions) ([]string, error) {
	var u SevenTVUser
	var err error

	if opt == nil {
		return nil, errors.New("No 7TV options provided.")
	} else if opt.SevenTVID != "" {
		u, err = get7TVUser(opt.SevenTVID)
		if err != nil {
			return nil, err
		}
	} else if opt.Platform != "" && opt.PlatformID != "" {
		u, err = get7TVPlatformUser(opt.Platform, opt.PlatformID)
		if err != nil {
			return nil, err
		}
		opt.SevenTVID = u.ID
	} else {
		return nil, errors.New("Either provide a 7TV ID or both platform and platform ID.")
	}
	ids := make([]string, 0, len(u.EmoteSets))
	for _, s := range u.EmoteSets {
		ids = append(ids, s.ID)
	}
	return ids, err
}
