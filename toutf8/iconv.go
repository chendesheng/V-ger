//
// iconv.go
//
package toutf8

// #cgo darwin  LDFLAGS: -liconv
// #cgo freebsd LDFLAGS: -liconv
// #cgo windows LDFLAGS: -liconv
// #include <iconv.h>
// #include <stdlib.h>
// #include <errno.h>
//
// size_t iconv_all(iconv_t cd, char *in, size_t inbytes, char **out, size_t bufsize) {
//   *out = (char *)malloc(bufsize);
//   char *outcurrent = *out;
//   size_t outbytes = 0;
//   size_t outbytesleft = bufsize;
//   while (inbytes > 0) {
//     iconv(cd, &in, &inbytes, &outcurrent, &outbytesleft);
//     outbytes += bufsize-outbytesleft;
//     *out = (char *)realloc(*out, outbytes+bufsize);
//     outcurrent = *out + outbytes;
//     outbytesleft = bufsize;
//   }
//   return outbytes;
// }
import "C"

import (
	"syscall"
	"unsafe"
)

var EILSEQ = syscall.Errno(C.EILSEQ)
var E2BIG = syscall.Errno(C.E2BIG)

const DefaultBufSize = 4096

type Iconv struct {
	Handle C.iconv_t
}

// Open returns a conversion descriptor cd, cd contains a conversion state and can not be used in multiple threads simultaneously.
func Open(tocode string, fromcode string) (cd Iconv, err error) {

	tocode1 := C.CString(tocode)
	defer C.free(unsafe.Pointer(tocode1))

	fromcode1 := C.CString(fromcode)
	defer C.free(unsafe.Pointer(fromcode1))

	ret, err := C.iconv_open(tocode1, fromcode1)
	if err != nil {
		return
	}
	cd = Iconv{ret}
	return
}

func (cd Iconv) Close() error {
	_, err := C.iconv_close(cd.Handle)
	return err
}

func (cd Iconv) ConvStr(s string) (string, error) {
  cs := C.CString(s)
  defer C.free(unsafe.Pointer(cs))

  out := (*C.char)(unsafe.Pointer(uintptr(0)))
  outbytes, errno := C.iconv_all(cd.Handle, cs, C.size_t(len(s)), (**C.char)(&out), DefaultBufSize);
  if (errno != nil) {
    return "", errno
  }
  return C.GoStringN(out, C.int(outbytes)), nil
}

