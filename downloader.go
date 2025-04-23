// TODO:
// - Get URL from Name
// - Get channel custom emotes
// - Get custom badges
// - Option: Periodic update checks? (go CheckUpdates)

// Reference:
// https://github.com/SevenTV/chatterino7/blob/chatterino7/src/providers/seventv/SeventvAPI.cpp

package emodl

import (
	"errors"
	"fmt"
	"sync"
	"text/template"
	"time"
)

var (
	apiPathTmpl, _   = template.New("api").Parse("/{{ .Version }}/{{ .Path }}")
	emotePathTmpl, _ = template.New("emote").Parse("emote/{{ .ID }}/{{ .Scale }}.{{ .Ext }}")

	imageFallbacks = [11]string{"WEBP", "AVIF", "APNG", "GIF", "PNG", "JPEG", "JPG", "JFIF", "PJPEG", "PJP", "SVG"}
)

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	ID     string `json:"id"`
}

type Emote struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Locations []string `json:"locations"`
	Images    []Image  `json:"images"`
}

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
	SevenTV *SevenTVOptions
}

// Downloads and caches third party emote data as maps indexed by name. 
type Downloader struct {
	Options       DownloaderOptions
	BTTVEmotes    map[string]BTTVEmote
	FFZEmotes     map[string]FFZEmote
	SevenTVEmotes map[string]SevenTVEmote
}

func NewDownloader(opt *DownloaderOptions) *Downloader {
	ed := &Downloader{Options: *opt}
	ed.BTTVEmotes = make(map[string]BTTVEmote, 64)
	ed.FFZEmotes = make(map[string]FFZEmote, 64)
	ed.SevenTVEmotes = make(map[string]SevenTVEmote, 64)
	return ed
}

// Loads all emote and badge data into memory based on configuration.
// Returns a map of emotes indexed by name.
func (ed *Downloader) Load() (map[string]Emote, error) {
	if ed == nil {
		return nil, errors.New("Nil dereference on Downloader")
	}
	var err error

	emotes := make(map[string]Emote, 256)

	errorChan := make(chan error, 8)
	sevenTVEmotesChan := make(chan SevenTVEmoteSet, 8)
	bttvEmotesChan := make(chan BTTVEmoteSlice, 8)

	wgdone := make(chan struct{})
	done := make(chan struct{})

	// Copier goroutine will copy emote data into the map as it comes in
	go func() {
		for {
			select {
			case set := <-sevenTVEmotesChan:
				for _, data := range set.Emotes {
					e := data.Data
					ed.SevenTVEmotes[e.Name] = e
					emotes[e.Name], err = e.AsEmote()
					if err != nil {
						errorChan <- err
					}
				}
			case es := <-bttvEmotesChan:
				for _, e := range es {
					ed.BTTVEmotes[e.Name] = e
					emotes[e.Name] = e.AsEmote()
				}
			case e := <-errorChan:
				err = errors.Join(e)
			case <-wgdone:
				// Collect all buffered errors and emotes when done downloading
				close(sevenTVEmotesChan)
				for s := range sevenTVEmotesChan {
					for _, data := range s.Emotes {
						e := data.Data
						ed.SevenTVEmotes[e.Name] = e
						emotes[e.Name], err = e.AsEmote()
						if err != nil {
							errorChan <- err
						}
					}
				}
				close(bttvEmotesChan)
				for es := range bttvEmotesChan {
					for _, e := range es {
						ed.BTTVEmotes[e.Name] = e
						emotes[e.Name] = e.AsEmote()
					}
				}
				close(errorChan)
				for e := range errorChan {
					err = errors.Join(e)
				}
				done <- struct{}{}
				return
			}
		}
	}()

	// Get request routines will download emote data asynchronously
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		bttvEmotes, err := getGlobalBTTVEmotes()
		if err != nil {
			errorChan <- errors.New(fmt.Sprintf("emodl: %v: failure getting global BTTV emotes", err))
			return
		}
		bttvEmotesChan <- bttvEmotes
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s, err := get7TVEmoteSet("global")
		if err != nil {
			errorChan <- errors.New(fmt.Sprintf("emodl: %v: failure getting global 7TV emotes", err))
			return
		}
		sevenTVEmotesChan <- s
	}()

	if ed.Options.SevenTV != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sids, err := get7TVUserEmoteSetIDs(ed.Options.SevenTV)
			if err != nil {
				errorChan <- errors.New(fmt.Sprintf("emodl: %v: failure getting user 7TV emotes with opt %v", err, *ed.Options.SevenTV))
				return
			}

			for _, sid := range sids {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s, err := get7TVEmoteSet(sid)
					if err != nil {
						errorChan <- errors.New(fmt.Sprintf("emodl: %v: failure getting 7TV emote set %s", err, sid))
						return
					}
					sevenTVEmotesChan <- s
				}()
			}
		}()
	}

	wg.Wait()
	wgdone <- struct{}{}

	timeout := time.NewTicker(5 * time.Second)
	select {
	case <-done:
	case <-timeout.C:
	}

	return emotes, err
}
