package audio

import (
	"fmt"
	// "io"
	"log"
	. "player/libav"
	"player/sdl"
)

type audioDriver interface{}

type sdlAudio struct {
	audioSpec sdl.AudioSpec
	volume    byte
}

func init() {
	if sdl.Init(sdl.SDL_INIT_AUDIO) {
		log.Print(sdl.GetError())
	}
}

func (a *sdlAudio) Open(channels int, sampleRate int, callback func(int) []byte) error {
	layout := GetChannelLayout("stereo")
	if channels == 1 {
		layout = GetChannelLayout("mono")
	}

	println("channels:", uint8(GetChannelLayoutNbChannels(layout)))

	var desired sdl.AudioSpec
	desired.Init()
	desired.SetFreq(sampleRate)
	desired.SetFormat(sdl.AUDIO_S16LSB)
	desired.SetChannels(uint8(GetChannelLayoutNbChannels(layout)))
	desired.SetSilence(0)
	desired.SetSamples(4096) //audio buffer size

	desired.SetCallback(func(userdata sdl.Object, stream sdl.Object, length int) {
		stream.SetZero(length)

		for length > 0 {
			p := callback(length)
			if len(p) > 0 {
				sdl.MixAudioFormat(&stream, p, a.audioSpec.Format(), int(float64(a.volume)/100*sdl.MIX_MAXVOLUME))
				// stream.Write(p)
				length -= len(p)
			}
		}
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
	a.volume-- //may overflow

	if a.volume < 0 || a.volume > 100 {
		a.volume = 0
	}
	return a.volume
}
