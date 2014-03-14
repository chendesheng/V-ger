package cld

//#cgo LDFLAGS: -L/Users/Roy/Vgerproj/src/cld/internal -lcld2 -lstdc++
//#cgo CFLAGS: -Ipublic
//#include "cld.h"
//#include <stdlib.h>
import "C"
import (
	// "io"
	"strings"
	"unsafe"
)

func DetectLanguage(s string) []string {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	res := (*_Ctype_char)(C.malloc(20))
	defer C.free(unsafe.Pointer(res))

	C.detect_language(cs, C.int(len(s)), res)
	return strings.Split(strings.TrimSpace(C.GoString(res)), ",")
}

func DetectLanguage2(s string) (string, string) {
	// println(s)
	langs := DetectLanguage(s)
	if len(langs) == 1 {
		langs = append(langs, "")
	}
	return langs[0], langs[1]
}
