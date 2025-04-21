// TODO:
// - Create interface for all emotes and badges
// - Get URL from Name
// - Option: Periodic update checks? (go CheckUpdates)

// Reference:
// https://github.com/SevenTV/chatterino7/blob/chatterino7/src/providers/seventv/SeventvAPI.cpp

package emodl

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

type DownloaderOptions struct {
	BTTV    bool
	FFZ     bool
	SevenTV bool
}

type Downloader struct {
	BTTVEmotes    map[string]BTTVEmote
	FFZEmotes     map[string]FFZEmote
	SevenTVEmotes map[string]SevenTVEmote
}

func NewDownloader(opt *DownloaderOptions) *Downloader {
	ed := &Downloader{}
	opt.FFZ = false // FFZ not implemented yet!
	if opt.BTTV {
		ed.BTTVEmotes = make(map[string]BTTVEmote, 64)
	}
	if opt.FFZ {
		ed.FFZEmotes = make(map[string]FFZEmote, 64)
	}
	if opt.SevenTV {
		ed.SevenTVEmotes = make(map[string]SevenTVEmote, 64)
	}
	return ed
}

// Loads all emote and badge data into memory based on configuration.
func (ed *Downloader) Load() error {
	if ed == nil {
		return errors.New("Nil dereference on Downloader")
	}
	var err error

	errorChan := make(chan error, 4)

	var wg sync.WaitGroup

	if ed.BTTVEmotes != nil {
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
	}

	if ed.SevenTVEmotes != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c, err := get7TVEmoteCollection("global")
			if err != nil {
				errorChan <- err
				return
			}

			for _, data := range c.Emotes {
				e := data.Data
				ed.SevenTVEmotes[e.Name] = e
			}
		}()
	}

	wg.Wait()

	close(errorChan)

	for e := range errorChan {
		err = errors.Join(e)
	}

	return err
}
