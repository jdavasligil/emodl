package emotedownloader

import "testing"

func TestGet7TVEmoteCollection(t *testing.T) {
	t.Parallel()
	c, err := get7TVEmoteCollection("global")
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("Nil collection")
	}
}
