// TODO:
// - Create interface for all emotes and badges
// - Get URL from Name

// Reference:
// https://github.com/SevenTV/chatterino7/blob/chatterino7/src/providers/seventv/SeventvAPI.cpp

package emotedownloader

import (
	"errors"
	"sync"
	"text/template"
)

var (
	apiPathTmpl, _   = template.New("api").Parse("/{{ .Version }}/{{ .Path }}")
	emotePathTmpl, _ = template.New("emote").Parse("emote/{{ .ID }}/{{ .Scale }}.{{ .Ext }}")
)

type emotePath struct {
	ID    string
	Scale string
	Ext   string
}

type apiPath struct {
	Version string
	Path    string
}

type FFZEmote struct {
}

type EmoteDownloaderConfig struct {
	BTTV    bool
	FFZ     bool
	SevenTV bool
}

type EmoteDownloader struct {
	BTTVEmotes    map[string]BTTVEmote
	FFZEmotes     map[string]FFZEmote
	SevenTVEmotes map[string]SevenTVEmote
}

func NewEmoteDownloader(config *EmoteDownloaderConfig) *EmoteDownloader {
	ed := &EmoteDownloader{}
	if config.BTTV {
		ed.BTTVEmotes = make(map[string]BTTVEmote, 64)
	}
	if config.FFZ {
		ed.FFZEmotes = make(map[string]FFZEmote, 64)
	}
	if config.SevenTV {
		ed.SevenTVEmotes = make(map[string]SevenTVEmote, 64)
	}
	return ed
}

// Loads all emote and badge data into memory based on configuration.
func (ed *EmoteDownloader) Load() error {
	var err error

	errorChan := make(chan error, 4)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		bttvEmotes, err := getGlobalBTTVEmotes()
		if err != nil {
			errorChan <- err
			return
		}

		for _, e := range bttvEmotes {
			ed.BTTVEmotes[e.Name] = e
		}
	}()

	wg.Wait()

	close(errorChan)

	for e := range errorChan {
		err = errors.Join(e)
	}

	return err
}
