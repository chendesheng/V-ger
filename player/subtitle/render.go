package subtitle

import (
	. "vger/player/shared"
)

type SubRender interface {
	SendShowText(SubItem)
	SendHideText(int)
}
