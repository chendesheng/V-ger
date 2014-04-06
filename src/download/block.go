package download

type block struct {
	from, to int64
	data     []byte
}

func (b *block) reset(from, to int64) {
	b.from, b.to = from, to
	size := to - from

	if int64(cap(b.data)) < size {
		b.data = make([]byte, size)
	} else {
		b.data = b.data[:size]
	}
}
