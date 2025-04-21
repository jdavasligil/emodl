package emodl

import (
	"testing"
)

func TestDownloader(t *testing.T) {
	t.Parallel()
	ed := NewDownloader(&DownloaderConfig{
		BTTV:    true,
		SevenTV: true,
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
	//s, err := prettyPrint(ed.BTTVEmotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}
