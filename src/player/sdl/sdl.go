package sdl

/*
#include "sdl2/sdl.h"
*/
import "C"

func Init(flags uint) bool {
	return C.SDL_Init(C.Uint32(flags)) != 0
}

func GetError() string {
	return C.GoString(C.SDL_GetError())
}

func OpenAudio(spec AudioSpec) (int, AudioSpec) {
	var obtained C.SDL_AudioSpec
	ret := C.SDL_OpenAudio(spec.ptr, &obtained)
	if ret == 0 {
		return int(ret), AudioSpec{ptr: &obtained}
	} else {
		return int(ret), AudioSpec{}
	}
}

func PauseAudio(vol int) {
	C.SDL_PauseAudio(C.int(vol))
}

func QuitSubSystem(flags uint) {
	C.SDL_QuitSubSystem(C.Uint32(flags))
}

func CloseAudio() {
	C.SDL_CloseAudio()
}

func Quit() {
	C.SDL_Quit()
}

func MixAudioFormat(dst Object, src []byte, format int, volume int) {
	dst.SetZero(len(src))
	C.SDL_MixAudioFormat((*_Ctype_Uint8)(dst.ptr), (*_Ctype_Uint8)(&src[0]), C.SDL_AudioFormat(format), C.Uint32(len(src)), C.int(volume))
}

const MIX_MAXVOLUME = C.SDL_MIX_MAXVOLUME
