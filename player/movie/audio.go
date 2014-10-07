package movie

import (
	"log"
	"strings"
	. "vger/player/libav"
	. "vger/player/movie/audio"
	. "vger/player/shared"
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
func (m *Movie) setupAudio() error {
	ctx := m.ctx

	audioStreams := ctx.AudioStream()
	m.audioStreams = audioStreams

	log.Print("setupAudio:", len(audioStreams), m.p.Volume)

	if len(audioStreams) > 0 {
		selectedStream := getDefaultAudioStream(audioStreams, m.p.SoundStream)
		selected := selectedStream.Index()
		m.p.SoundStream = selected
		SavePlayingAsync(m.p)

		var err error
		m.a = NewAudio(m.c, m.p.Volume)

		err = m.a.Open(selectedStream)
		if err != nil {
			return err
		}
	}

	return nil

}
