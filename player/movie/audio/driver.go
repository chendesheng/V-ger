package audio

import (
	"log"
	"sync"
	"unsafe"

	"github.com/gordonklaus/portaudio"
)

func init() {
	err := portaudio.Initialize()
	if err != nil {
		log.Print(err)
	}
}

type audioDriver interface{}

type portAudio struct {
	sync.Mutex
	volume int
	stream *portaudio.Stream
}

func (a *portAudio) Open(sampleRate int, callback func(int) []byte) error {
	var err error
	a.stream, err = portaudio.OpenDefaultStream(0, 1, float64(sampleRate), portaudio.FramesPerBufferUnspecified, func(out []int32) {
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

	return a.stream.Start()
}

func (a *portAudio) Close() {
	if a.stream != nil {
		err := a.stream.Stop()
		if err != nil {
			log.Print(err)
		}
	}
}

func (a *portAudio) IncreaseVolume() int {
	a.Lock()
	defer a.Unlock()

	a.volume += 1

	if a.volume > 16 {
		a.volume = 16
	}
	return a.volume
}
func (a *portAudio) DecreaseVolume() int {
	a.Lock()
	defer a.Unlock()

	a.volume -= 1

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
	v := float64(a.volume) / 10
	v2 := v * v
	return v2 * 1.2
}
