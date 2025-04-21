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
	t.Run("TestEmoteGetImageURL", func(t *testing.T) {
		e := c.Emotes[0].Data
		s, err := prettyPrint(e.Host.Files)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
		url, err := e.GetImageURL("2x", "webp")
		t.Logf("\n\tURL: %s\n", url)
	})
}
