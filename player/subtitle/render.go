package subtitle

import (
	. "vger/player/shared"
)

type SubRender interface {
	SendShowText(SubItemArg)
	SendHideText(SubItemArg)
}
