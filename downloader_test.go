package emodl

import (
	"testing"
)

func TestDownloader(t *testing.T) {
	t.Parallel()
	ed := NewDownloader(&DownloaderOptions{
		BTTV: true,
		SevenTV: &SevenTVOptions{
			Platform:   "twitch",
			PlatformID: "1048391821",
		},
	})
	if ed == nil {
		t.Fatal("Nil downloader")
	}
	err := ed.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(ed.BTTVEmotes) == 0 {
		t.Fatal("BTTVEmotes is empty.")
	}
	if len(ed.SevenTVEmotes) == 0 {
		t.Fatal("SevenTVEmotes is empty")
	}
	t.Logf("\nDownloaded:\n\t%d Emotes from 7TV\n\t%d Emotes from BTTV\n", len(ed.SevenTVEmotes), len(ed.BTTVEmotes))
	//s, err := prettyPrint(ed.BTTVEmotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}
