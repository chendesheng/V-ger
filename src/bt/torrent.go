package bt

import (
	"io"
)

type torrentFile struct {
	length int64
	md5sum string
	path   []string
}
type torrentInfo struct {
	pieceLength int
	pieces      []byte
	private     int
	name        string
	files       []torrentFile
}
type torrent struct {
	info         map[string]torrentInfo
	announceURL  string
	announceList []string
	creationDate int
	comment      string
	createdBy    string
	encoding     string
}

func parse(r io.Reader) torrent {

}
