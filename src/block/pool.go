package block

import "sync"

type BlockPool struct {
	sync.Pool
}

func (bp *BlockPool) Get(from int64, size int) *Block {
	bk := bp.Pool.Get().(*Block)
	bk.Reset(from, size)
	return bk
}

func (bp *BlockPool) Put(bk *Block) {
	bp.Pool.Put(bk)
}

func (bp *BlockPool) GetBlocks(count int, size int) []*Block {
	blocks := make([]*Block, 0, count)
	for i := 0; i < count; i++ {
		blocks = append(blocks, DefaultBlockPool.Get(0, size))
	}
	return blocks
}
func (bp *BlockPool) PutBlocks(blocks []*Block) {
	for _, bk := range blocks {
		DefaultBlockPool.Put(bk)
	}
}

var DefaultBlockPool BlockPool

func init() {
	DefaultBlockPool.New = func() interface{} {
		return &Block{0, nil}
	}
}
