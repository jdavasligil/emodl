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
		for batch := range imgIter {
			if ed.Err != nil {
				t.Fatal(ed.Err)
			}
			if batch == nil {
				t.Fatal("Nil batch.")
			}
			if len(batch) == 0 {
				t.Fatal("Batch length is 0.")
			}
			if batch[0].Image == nil {
				t.Fatal("Nil image.")
			}
			if batch[0].ID == "" {
				t.Fatal("Empty ID.")
			}
			t.Logf("ID: %s\tH: %d\tW: %d", batch[0].ID, batch[0].Image.Bounds().Dy(), batch[0].Image.Bounds().Dx())
		}
		if ed.Err != nil {
			t.Fatal(ed.Err)
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
