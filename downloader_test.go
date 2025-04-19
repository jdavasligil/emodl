package emotedownloader

import (
	"encoding/json"
	"testing"
)

func TestDownloader(t *testing.T) {
	ed := NewEmoteDownloader(&EmoteDownloaderConfig{
		BTTV: true,
	})
	err := ed.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(ed.BTTVEmotes) == 0 {
		t.Fatal("Emotes is empty.")
	}
	s, err := prettyPrint(ed.BTTVEmotes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func prettyPrint(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
