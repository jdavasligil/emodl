package emotedownloader

import (
	"encoding/json"
	"testing"
)

func TestDownloader(t *testing.T) {
	ed := &EmoteDownloader{}
	err := ed.GetBTTVGlobalEmotes()
	if err != nil {
		t.Fatal(err)
	}
	s, err := PrettyPrint(ed.Emotes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func PrettyPrint(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
