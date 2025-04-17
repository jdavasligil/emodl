package emotedownloader

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestDownloader(t *testing.T) {
	ed := &EmoteDownloader{}
	err := ed.Download()
	if err != nil {
		t.Fatal(err)
	}
	if len(ed.Emotes) == 0 {
		t.Fatal("Emotes is empty.")
	}
	if !slices.Contains(EmoteProviders, ed.Emotes[0].Provider) {
		t.Fatal("Provider not recognized.")
	}
	//s, err := prettyPrint(ed.Emotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}

func prettyPrint(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
