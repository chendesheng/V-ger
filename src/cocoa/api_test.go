package cocoa

import (
	"testing"
)

func TestTrash(t *testing.T) {
	api := cocoaNativeAPI{}
	api.MoveFileToTrash("/Volumes/Data/Downloads/Video", "vger copy.db")
}
