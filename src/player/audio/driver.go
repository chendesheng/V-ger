package audio

import (
	"sync"
	"unsafe"

	"code.google.com/p/portaudio-go/portaudio"
)

func init() {
	portaudio.Initialize()
}

type audioDriver interface{}

type portAudio struct {
	sync.Mutex
	volume float64
	stream *portaudio.Stream
}

func (a *portAudio) Open(channels int, sampleRate int, callback func(int) []byte) error {

	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	args := portaudio.HighLatencyParameters(nil, h.DefaultOutputDevice)
	args.SampleRate = float64(sampleRate)
	args.Output.Channels = 1

	a.stream, err = portaudio.OpenStream(args, func(out []int32) {
		// println(out)
		length := len(out)
		for length > 0 {
			p := callback(length * 4)
			data := (*(*[]int32)(unsafe.Pointer(&p)))[:len(p)/4]
			if len(data) > 0 {
				off := len(out) - length
				for i, b := range data {
					out[off+i] = int32(float64(b)*a.getVolume() + 0.5)
				}

				length -= len(data)
			}
		}
	})
	if err != nil {
		return err
	}

	a.stream.Start()

	return nil
}

func (a *portAudio) Close() {
	a.stream.Stop()
	a.stream.Close()
}

func (a *portAudio) IncreaseVolume() float64 {
	a.Lock()
	defer a.Unlock()

	a.volume += 0.01

	if a.volume > 1 {
		a.volume = 1
	}
	return a.volume
}
func (a *portAudio) DecreaseVolume() float64 {
	a.Lock()
	defer a.Unlock()

	a.volume -= 0.01

	if a.volume < 0 {
		a.volume = 0
	}
	return a.volume
}
func (a *portAudio) getVolume() float64 {
	a.Lock()
	defer a.Unlock()

	//linear volume
	//check this: http://www.dr-lex.be/info-stuff/volumecontrols.html
	v2 := a.volume * a.volume
	return v2 * v2 * 1.2
}
