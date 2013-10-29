#include "_cgo_export.h"

void callback(void *userdata, Uint8 * stream, int len) {
	goCallback(userdata, stream, len);
}
void SetCallbackCB(SDL_AudioSpec *spec) {
	spec->callback = callback;
}

