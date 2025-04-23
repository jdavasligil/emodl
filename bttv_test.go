package emodl

import (
	"testing"
)

func TestGetGlobalBTTVEmotes(t *testing.T) {
	t.Parallel()
	bttvEmotes, err := getGlobalBTTVEmotes()
	if err != nil {
		t.Fatal(err)
	}
	if bttvEmotes == nil {
		t.Fatal("nil emotes slice")
	}

	t.Logf("BTTVEmoteSlice Size: %s", humanSize(bttvEmotes.Size()))

	for _, e := range bttvEmotes {
		if e.ID == "" {
			t.Logf("ID empty: %v", e)
			t.Fail()
		}
		if e.Name == "" {
			t.Logf("Name empty: %v", e)
			t.Fail()
		}
		if e.ImageType == "" {
			t.Logf("ImageType empty: %v", e)
			t.Fail()
		}
	}
	//s, err := prettyPrint(bttvEmotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}
