package util

import (
	"testing"
)

func TestCleanMovieName(t *testing.T) {
	if CleanMovieName("Carlito's.Way.1993.720p.BluRay.DTS.x264-dxva.mkv.vger-task.txt") == "Carlito's Way 1993 720p" {
		t.Error("name not equal")
	}

	if CleanMovieName("The.Devils.Advocate.1997.UNRATED.DC.1080p.BluRay.X264-AMIABLE.mkv.vger-task.txt") == "The Devils Advocate 1997 UNRATED DC" {
		t.Error("name not equal")
	}
	if CleanMovieName("Google IO 2013 - Advanced Go Concurrency Patterns [720p].mp4.vger-task.txt") == "Google IO 2013" {
		t.Error("name not equal")
	}
	if CleanMovieName("breaking.bad.s01e06.720p.bluray.x264-reward.mkv.vger-task.txt") == "breaking bad s01e06" {
		t.Error("name not equal")
	}
	if CleanMovieName("House.of.Cards.S01E07.WEBRip.720p.H.264.AAC.2.0-HoC.mkv.vger-task.txt") == "House of Cards S01E07" {
		t.Error("name not equal")
	}
	if CleanMovieName("Scarface.1983.1080p.BluRay.X264-AMIABLE.mkv.vger-task.txt") == "Scarface 1983" {
		t.Error("name not equal")
	}
}
