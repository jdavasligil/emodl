package emodl

import (
	"testing"
)

func TestDownloader(t *testing.T) {
	t.Parallel()
	ed := NewDownloader(DownloaderOptions{
		BTTV: &BTTVOptions{
			Platform:   "twitch",
			PlatformID: "39226538",
		},
		SevenTV: &SevenTVOptions{
			Platform:   "twitch",
			PlatformID: "1048391821",
		},
	})
	emotes, err := ed.Load()
	if err != nil {
		t.Fatal(err)
	}
	if emotes == nil {
		t.Fatal("Emotes slice is nil")
	}
	if len(ed.BTTVEmotes) == 0 {
		t.Fatal("BTTVEmotes is empty.")
	}
	if len(ed.SevenTVEmotes) == 0 {
		t.Fatal("SevenTVEmotes is empty")
	}
	if len(emotes) != (len(ed.BTTVEmotes) + len(ed.SevenTVEmotes)) {
		if len(ed.BTTVEmotes) < len(ed.SevenTVEmotes) {
			for name := range ed.BTTVEmotes {
				_, ok := ed.SevenTVEmotes[name]
				if ok {
					t.Logf("Emote Conflict: [BTTV] [7TV] -> %s", name)
				}
			}
		} else {
			for name := range ed.SevenTVEmotes {
				_, ok := ed.BTTVEmotes[name]
				if ok {
					t.Logf("Emote Conflict: [BTTV] [7TV] -> %s", name)
				}
			}

		}
	}
	t.Logf("\nDownloaded:\n\t%d Emotes from 7TV\n\t%d Emotes from BTTV\n\t%d Emotes Total\n", len(ed.SevenTVEmotes), len(ed.BTTVEmotes), len(emotes))
	//s, err := prettyPrint(emotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}
