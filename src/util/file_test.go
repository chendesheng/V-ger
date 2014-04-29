package util

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	os.RemoveAll("95edab554826f0e0ebdd0205a3f94dbf")

	resp, err := http.Get("http://res.yyets.com/ftp/2014/0220/95edab554826f0e0ebdd0205a3f94dbf.rar")
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	ioutil.WriteFile("./95edab554826f0e0ebdd0205a3f94dbf.rar", data, 0666)

	Extract("./unar", "./95edab554826f0e0ebdd0205a3f94dbf.rar")

	_, err = os.Stat("95edab554826f0e0ebdd0205a3f94dbf/house.of.cards.2013.s02e11.720p.nf.webrip.dd5.1.x264-ntb/House.of.Cards.2013.S02E11.720p.NF.WEBRip.DD5.1.x264-NTb.繁体.ass")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEmulateFiles(t *testing.T) {
	files := make([]string, 0)
	EmulateFiles("./testfolder", func(filename string) {
		files = append(files, filename)
	}, "txt")

	if len(files) != 1 {
		t.Errorf("Expect 1 files but %d", len(files))
	}

	println(files[0])
}
