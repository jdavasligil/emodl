package emodl

import (
	"fmt"
	"strings"
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
		FFZ: &FFZOptions{
			Platform:   "twitch",
			PlatformID: "39226538",
		},
	})

	emotes, err := ed.Load()
	if err != nil {
		t.Fatal(err)
	}
	if emotes == nil {
		t.Fatal("Emotes map is nil")
	}
	if len(ed.BTTVEmotes) == 0 {
		t.Fatal("BTTVEmotes is empty.")
	}
	if len(ed.SevenTVEmotes) == 0 {
		t.Fatal("SevenTVEmotes is empty")
	}
	if len(ed.FFZEmotes) == 0 {
		t.Fatal("FFZEmotes is empty")
	}

	t.Log(ed.ReportConflicts(emotes))

	var sb strings.Builder
	sb.WriteString("\nDownloaded:\n")
	sb.WriteString(fmt.Sprintf("\t7TV:   %d\n", len(ed.SevenTVEmotes)))
	sb.WriteString(fmt.Sprintf("\tBTTV:  %d\n", len(ed.BTTVEmotes)))
	sb.WriteString(fmt.Sprintf("\tFFZ:   %d\n", len(ed.FFZEmotes)))
	sb.WriteString(fmt.Sprintf("\tTotal: %d\n", len(ed.FFZEmotes)+len(ed.SevenTVEmotes)+len(ed.BTTVEmotes)))
	t.Log(sb.String())
	//s, err := prettyPrint(emotes)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(s)
}
