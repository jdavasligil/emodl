package emodl

import (
	"bytes"
	"testing"

	"github.com/mailru/easyjson"
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
			t.Fatal("ID empty")
		}
		if e.Name == "" {
			t.Fatal("Name empty")
		}
		if e.ImageType == "" {
			t.Fatal("ImageType empty")
		}
		if e.UserID == "" {
			t.Fatal("UserID empty")
		}
	}
}

// TODO: JSON fuzzing with arbitrary possible json values
func TestBTTVEasyJSON(t *testing.T) {
	t.Parallel()
	type BTTVJSONTest struct {
		EmotesBefore BTTVEmoteSlice
		EmotesAfter  BTTVEmoteSlice
	}
	tests := []BTTVJSONTest{
		{BTTVEmoteSlice{{}}, BTTVEmoteSlice{{}}},
		{BTTVEmoteSlice{{ID: "1", Name: "1", ImageType: "png", Animated: true, UserID: "1"}}, BTTVEmoteSlice{{}}},
	}

	w := &bytes.Buffer{}
	for _, test := range tests {
		easyjson.MarshalToWriter(test.EmotesBefore, w)
		easyjson.UnmarshalFromReader(w, &test.EmotesAfter)
		for i := range test.EmotesBefore {
			if test.EmotesAfter[i] != test.EmotesBefore[i] {
				t.Fatalf("Expected: %v\n Got: %v\n", test.EmotesBefore, test.EmotesAfter)
			}
		}
		w.Reset()
	}
}
