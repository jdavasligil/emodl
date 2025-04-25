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
	s, err := prettyPrint(sets)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}
