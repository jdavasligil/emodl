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
	if len(ed.BTTVEmotes) == 0 {
		t.Fatal("Emotes is empty.")
	}
	if !slices.Contains(EmoteProviders, ed.BTTVEmotes[0].Provider) {
		t.Fatal("Provider not recognized.")
	}
	s, err := prettyPrint(ed.BTTVEmotes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	t.Run("TestBTTVImages", func(t *testing.T) {
		imgIter := ed.BTTVImages(ImageScale1X)
		for id, img := range imgIter {
			if ed.Err != nil {
				t.Fatal(ed.Err)
			}
			if img == nil {
				t.Fatal("Nil image batch.")
			}
			if len(img) == 0 {
				t.Fatal("No image batch.")
			}
			batch := img[0]
			if batch == nil {
				t.Fatal("Nil batch")
			}
			t.Logf("ID: %s\tH: %d\tW: %d", id, batch.Bounds().Dy(), batch.Bounds().Dx())
		}
	})
}

func prettyPrint(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
