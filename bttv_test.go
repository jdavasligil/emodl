package emodl

import (
	"testing"
)

func TestGetGlobalBTTVEmotes(t *testing.T) {
	t.Parallel()
	bttvEmotes, err := getBTTVGlobalEmotes()
	if err != nil {
		t.Fatal(err)
	}
	if bttvEmotes == nil {
		t.Fatal("Emote slice should never be nil")
	}
	if len(bttvEmotes) == 0 {
		t.Fatal("No emotes obtained.")
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
	}
	//s, err := prettyPrint(bttvEmotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}

func TestGetBTTVUserEmotes(t *testing.T) {
	t.Parallel()
	bttvEmotes, err := getBTTVUserEmotes("twitch", "39226538")
	if err != nil {
		t.Fatal(err)
	}
	if bttvEmotes == nil {
		t.Fatal("bttv emotes can not be nil")
	}
	if len(bttvEmotes) == 0 {
		t.Fatal("No emotes obtained")
	}
}
