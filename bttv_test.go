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
