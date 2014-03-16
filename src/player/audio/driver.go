package audio

import (
	"fmt"
	"io"
	"log"
	. "player/libav"
	"player/sdl"
)

type audioDriver interface{}

type sdlAudio struct {
	audioSpec sdl.AudioSpec
	volume    byte

	slience []byte

	stream sdl.Object
}

func (a *sdlAudio) Write(p []byte) (int, error) {
	sdl.MixAudioFormat(a.stream, p, a.audioSpec.Format(), int(float64(a.volume)/100*sdl.MIX_MAXVOLUME))
	a.stream.Offset(len(p))
	return len(p), nil
}

func init() {
	if sdl.Init(sdl.SDL_INIT_AUDIO) {
		log.Print(sdl.GetError())
	}
}

func (a *sdlAudio) Open(channels int, sampleRate int, callback func(io.Writer, int)) error {
	layout := GetChannelLayout("stereo")
	if channels == 1 {
		layout = GetChannelLayout("mono")
	}

	var desired sdl.AudioSpec
	desired.Init()
	desired.SetFreq(sampleRate)
	desired.SetFormat(sdl.AUDIO_S16LSB)
	desired.SetChannels(uint8(GetChannelLayoutNbChannels(layout)))
	desired.SetSilence(0)
	desired.SetSamples(4096) //audio buffer size

	desired.SetCallback(func(userdata sdl.Object, stream sdl.Object, length int) {
		a.stream = stream
		callback(a, length)
	})

	res, obtained := sdl.OpenAudio(desired)
	if res < 0 {
		return fmt.Errorf("sdl open audio error: ", sdl.GetError())
	}
	if obtained.IsNil() {
		return fmt.Errorf("sdl get nil obtained audio spec")
	}
	a.audioSpec = obtained

	sdl.PauseAudio(0)
	return nil
}

func (a *sdlAudio) Close() {
	sdl.CloseAudio()
}

func (a *sdlAudio) IncreaseVolume() byte {
	a.volume++

	if a.volume > 100 {
		a.volume = 100
	}
	return a.volume
}
func (a *sdlAudio) DecreaseVolume() byte {
	a.volume--

	if a.volume < 0 {
		a.volume = 0
	}
	return a.volume
}
