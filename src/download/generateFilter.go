package download

import (
	"block"
	"log"
)

type generateFilter struct {
	basicFilter
	from        int64
	to          int64
	blockSize   int
	chBlockSize chan int
	// blocks      []*block.Block
	maxConnections int
}

func (gf *generateFilter) nextBlock() (*block.Block, bool) {
	from := gf.from
	blockSize := gf.blockSize
	if from+int64(blockSize) > gf.to {
		blockSize = int(gf.to - from)
	}

	gf.from = from + int64(blockSize)

	if blockSize == 0 {
		gf.closeOutput()
		gf.drainInput()
		return nil, false
	}
	return block.DefaultBlockPool.Get(from, blockSize), true
}

func (gf *generateFilter) active() {
	if gf.blockSize == 0 {
		//small blocksize for fast boot
		gf.blockSize = 128 * block.KB
	}

	log.Printf("generateFilter:%d %d %d %d\n", gf.from, gf.to, gf.blockSize, gf.maxConnections)

	//boot
	for i := 0; i < gf.maxConnections; i++ {
		if bk, ok := gf.nextBlock(); ok {
			// println("generateFilter:", bk.From, len(bk.Data))
			if gf.writeOutput(*bk) {
				break
			}
		} else {
			return
		}
	}

	//change to a larger blocksize after boot
	gf.blockSize = 512 * block.KB
	for {
		select {
		case _, ok := <-gf.input:
			if !ok {
				return
			} else {
				if bk, ok := gf.nextBlock(); ok {
					gf.writeOutput(*bk)
				} else {
					return
				}
			}
			break
		case maxSpeed := <-gf.chBlockSize:
			log.Print("set block size: ", maxSpeed)
			gf.blockSize = getBlockSize(maxSpeed)
			break
		case <-gf.quit:
			gf.closeOutput()
			log.Print("quit generate block")
			return
		}
	}
}

func getBlockSize(maxSpeed int) int {
	if maxSpeed > 0 {
		return maxSpeed * block.KB
	} else {
		return 512 * block.KB
	}
}
