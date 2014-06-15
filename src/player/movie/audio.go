package movie

import (
	"fmt"
	"log"
	. "player/audio"
	. "player/gui"
	. "player/libav"
	. "player/shared"
	"strings"
)

func getStream(streams []AVStream, index int) AVStream {
	for _, stream := range streams {
		if stream.Index() == index {
			return stream
		}
	}

	return AVStream{}
}
func getStreamByLanguage(streams []AVStream, lang string) AVStream {
	for _, stream := range streams {
		dic := stream.MetaData()
		mp := dic.Map()
		l := strings.ToLower(mp["language"])
		if strings.Contains(l, lang) {
			return stream
		}
	}

	return AVStream{}
}
func getDefaultAudioStream(streams []AVStream, lastSelected int) AVStream {
	var selectedStream AVStream
	if selectedStream = getStream(streams, lastSelected); selectedStream.IsNil() {
		if selectedStream = getStreamByLanguage(streams, "eng"); selectedStream.IsNil() {
			selectedStream = streams[0]
		}
	}

	return selectedStream
}
func (m *Movie) setupAudioMenu(selected int) {
	HideAudioMenu()
	audioStreams := m.audioStreams
	audioStreamNames := make([]string, 0)
	audioStreamIndexes := make([]int32, 0)
	if len(audioStreams) > 1 {
		for _, stream := range audioStreams {
			dic := stream.MetaData()
			mp := dic.Map()
			title := mp["title"]                        //dic.AVDictGet("title", AVDictionaryEntry{}, 2).Value()
			language := strings.ToLower(mp["language"]) //dic.AVDictGet("language", AVDictionaryEntry{}, 2).Value()

			// println(title, language)
			audioStreamNames = append(audioStreamNames, fmt.Sprintf("[%s] %s", language, title))
			audioStreamIndexes = append(audioStreamIndexes, int32(stream.Index()))
		}
		m.w.InitAudioMenu(audioStreamNames, audioStreamIndexes, selected)
	}
}
func (m *Movie) setupAudio() {
	ctx := m.ctx

	audioStreams := ctx.AudioStream()
	m.audioStreams = audioStreams

	log.Print("setupAudio:", len(audioStreams))

	if len(audioStreams) > 0 {
		selectedStream := getDefaultAudioStream(audioStreams, m.p.SoundStream)
		selected := selectedStream.Index()
		m.p.SoundStream = selected
		SavePlayingAsync(m.p)

		var err error
		m.a = NewAudio(m.c, float64(m.p.Volume)/100)

		err = m.a.Open(selectedStream)
		if err != nil {
			log.Print(err)
			return
		}

		m.setupAudioMenu(selected)
	}

}
