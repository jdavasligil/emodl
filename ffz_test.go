package emodl

import "testing"

func TestGetFFZEmoteSet(t *testing.T) {
	sets, err := getFFZEmoteSets("global")
	if err != nil {
		t.Fatal(err)
	}
	if sets == nil {
		t.Log("Sets must not be nil.")
		t.Fail()
	}
	if len(sets) == 0 {
		t.Log("Sets length is zero")
		t.Fail()
	}
	// s, err := prettyPrint(sets)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(s)
}

func TestGetFFZRoomSet(t *testing.T) {
	set, err := getFFZRoomEmoteSet("twitch", "39226538")
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Emotes) == 0 {
		t.Fatal("Emote set is empty")
	}
	// s, err := prettyPrint(set)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(s)
}
