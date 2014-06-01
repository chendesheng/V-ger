package download

type block struct {
	from int64
	data []byte
}

func (b *block) reset(from int64, size int) {
	b.from = from

	if cap(b.data) < size {
		b.data = make([]byte, size)
	} else {
		b.data = b.data[:size]
	}
}
