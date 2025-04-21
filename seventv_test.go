package emodl

import (
	"testing"
)

func TestGet7TVEmoteCollection(t *testing.T) {
	t.Parallel()
	c, err := get7TVEmoteCollection("global")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Nil collection")
	}
	if len(c.Emotes) == 0 {
		t.Fatal("No 7TV emotes")
	}
	t.Logf("7TVEmoteCollection Size: %s", humanSize(c.Size()))
	t.Run("TestEmoteGetImage", func(t *testing.T) {
		e := c.Emotes[0].Data
		img, err := e.GetImage("2x", "webp")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("\nURL: %s\nID: %s\n", img.URL, img.ID)
		img, err = e.GetImage("", "")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("\nURL: %s\nID: %s\n", img.URL, img.ID)
	})
}
