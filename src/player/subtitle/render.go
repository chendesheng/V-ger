package subtitle

import (
	. "player/shared"
)

type SubRender interface {
	SendShowText(SubItemArg)
	SendHideText(SubItemArg)
}
