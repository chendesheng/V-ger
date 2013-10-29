package sdl

/*
#include "sdl2/sdl.h"
void SetCallbackCB(SDL_AudioSpec *spec);
*/
import "C"
import (
	"unsafe"
)

type AudioSpec struct {
	ptr *C.SDL_AudioSpec
}

func (spec *AudioSpec) Init() {
	spec.ptr = &C.SDL_AudioSpec{}
}

func (spec *AudioSpec) IsNil() bool {
	return spec.ptr == nil
}
func (spec *AudioSpec) SetFreq(freq int) {
	spec.ptr.freq = C.int(freq)
}
func (spec *AudioSpec) SetFormat(format int) {
	spec.ptr.format = C.SDL_AudioFormat(format)
}
func (spec *AudioSpec) SetChannels(channels uint8) {
	spec.ptr.channels = C.Uint8(channels)
}
func (spec *AudioSpec) SetSilence(silence int) {
	spec.ptr.silence = C.Uint8(silence)
}
func (spec *AudioSpec) SetSamples(samples int) { //audio buffer size
	spec.ptr.samples = C.Uint16(samples)
}

var callback func(userdata Object, stream Object, length int)

//export goCallback
func goCallback(userdata unsafe.Pointer, stream unsafe.Pointer, length int) {
	callback(Object{ptr: userdata}, Object{ptr: stream}, length)
}

func (spec *AudioSpec) SetCallback(fn func(userdata Object, stream Object, length int)) {
	callback = fn
	C.SetCallbackCB(spec.ptr)
}
