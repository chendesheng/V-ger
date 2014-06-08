package download

import "log"

// "time"
// "util"

type generateFilter struct {
	basicFilter
	from           int64
	to             int64
	blockSize      int64
	chBlockSize    chan int64
	maxConnections int
}

func (gf *generateFilter) nextBlock() (int64, int64, bool) {
	from := gf.from
	blockSize := gf.blockSize
	if from+blockSize > gf.to {
		blockSize = gf.to - from
	}

	gf.from = from + blockSize

	if blockSize == 0 {
		gf.closeOutput()
		gf.drainInput()
		return 0, 0, false
	}
	return from, blockSize, true
}

func (gf *generateFilter) active() {
	if gf.blockSize == 0 {
		//small blocksize for fast boot
		gf.blockSize = int64(128 * 1024)
	}

	//boot
	for i := 0; i < gf.maxConnections; i++ {
		if from, blockSize, ok := gf.nextBlock(); ok {
			if gf.writeOutput(block{from, make([]byte, blockSize, 512*1024)}) {
				return
			}
		}
	}

	//change to a larger blocksize after boot
	gf.blockSize = 512 * 1024
	for {
		select {
		case b, ok := <-gf.input:
			if !ok {
				return
			} else {
				if from, blockSize, ok := gf.nextBlock(); ok {
					b.reset(from, int(blockSize))
					gf.writeOutput(b)
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
			log.Print("quit generate block")
			return
		}
	}
}

func getBlockSize(maxSpeed int64) int64 {
	if maxSpeed > 0 {
		return maxSpeed * 1024
	} else {
		return 512 * 1024
	}
}
