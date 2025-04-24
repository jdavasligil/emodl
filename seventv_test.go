package emodl

import (
	"log"
	"testing"
)

func TestGet7TVEmoteSet(t *testing.T) {
	t.Parallel()
	s, err := get7TVEmoteSet("global")
	if err != nil {
		t.Fatal(err)
	}
	if s.Emotes == nil {
		t.Fatal("Nil set")
	}
	if len(s.Emotes) == 0 {
		t.Fatal("No 7TV emotes")
	}
	t.Logf("7TVEmoteSet Size: %s", humanSize(s.Size()))
	t.Run("TestEmoteGetImage", func(t *testing.T) {
		e := s.Emotes[0].Data
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

func TestGet7TVEmoteSetIDs(t *testing.T) {
	t.Parallel()
	t.Run("WithPlatformNoID", func(t *testing.T) {
		t.Parallel()
		sids, err := get7TVUserEmoteSetIDs("twitch", "")
		if err == nil {
			t.Logf("No error with no platform id")
			t.Fail()
		}
		if sids == nil {
			log.Fatal("Emote Set IDs should never be nil")
		}
		if len(sids) != 0 {
			t.Logf("Emote Set IDs not empty with no platform id")
			t.Fail()
		}
	})
	t.Run("WithPlatformIDNoPlatform", func(t *testing.T) {
		t.Parallel()
		sids, err := get7TVUserEmoteSetIDs("", "1048391821")
		if err == nil {
			t.Logf("No error with no platform")
			t.Fail()
		}
		if sids == nil {
			t.Logf("Emote Set IDs should never be nil")
			t.Fail()
		}
		if len(sids) != 0 {
			t.Logf("Emote Set IDs not empty with no platform")
			t.Fail()
		}
	})
	t.Run("WithPlatformAndPID", func(t *testing.T) {
		t.Parallel()
		sids, err := get7TVUserEmoteSetIDs("twitch", "1048391821")
		if err != nil {
			t.Fatal(err)
		}
		if len(sids) != 1 {
			t.Fatalf("Emote Set IDs: %v has more than one emote set", sids)
		}
		if sids[0] != "01JSG36904T5GM79JJXBVTSFKS" {
			t.Fatalf("Emote Set ID: %s different from expected %s", sids[0], "01JSG36904T5GM79JJXBVTSFKS")
		}
	})
}
