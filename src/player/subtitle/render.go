package subtitle

import (
	. "player/shared"
)

type SubRender interface {
	SendShowText(*SubItem) uintptr
	SendHideText(uintptr)
}
