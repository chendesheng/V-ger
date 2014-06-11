package block

type Block struct {
	From int64
	Data []byte
}

func (bk *Block) Reset(from int64, size int) {
	bk.From = from

	if cap(bk.Data) < size {
		bk.Data = make([]byte, size)
	} else {
		bk.Data = bk.Data[:size]
	}
}

func (bk *Block) Inside(position int64) bool {
	return bk.From <= position && position < bk.From+int64(len(bk.Data))
}
