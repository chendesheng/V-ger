package download

import (
	"fmt"
	"log"
	// "time"
	// "util"
)

type generateFilter struct {
	basicFilter
	from           int64
	to             int64
	blockSize      int64
	chBlockSize    chan int64
	maxConnections int
}

func (gf *generateFilter) active() {
	generateBlock(gf.input, gf.output, gf.chBlockSize, gf.from, gf.to, gf.blockSize, gf.maxConnections, gf.quit)
}

func generateBlock(input chan *block, output chan<- *block, chBlockSize chan int64, from, size int64, blockSize int64, maxConnections int, quit <-chan bool) {
	log.Printf("generate block output: %v", output)
	if blockSize == 0 {
		//small blocksize for fast boot
		blockSize = int64(32 * 1024)
	}

	to := from + blockSize
	if to > size {
		to = size
	}

	//boot
	for i := 0; i < maxConnections; i++ {
		select {
		case output <- &block{from, to, make([]byte, to-from)}:
			from = to
			to = from + blockSize
			if to > size {
				to = size
			}
			break
		case <-quit:
			return
		}
	}

	//change to a larger blocksize after boot
	blockSize = 512 * 1024

	log.Printf("output %v", output)
	maxSpeed := int64(0)
	for {
		select {
		case b, ok := <-input:
			if !ok {
				return
			} else {
				b.reset(from, to)
				select {
				case output <- b:
					if to == size {
						fmt.Println("return generate block ", size)
						close(output)
						for {
							select {
							case _, ok := <-input:
								if !ok {
									return
								}
							case <-quit:
								return
							}
						}
					} else {
						from = to
						to = from + blockSize
						if to > size {
							to = size
						}
					}
				case maxSpeed = <-chBlockSize:
					log.Print("set block size: ", maxSpeed)
					blockSize = getBlockSize(maxSpeed)
					break
				case <-quit:
					log.Print("quit generate block")
					return
				}
			}
			break
		case maxSpeed = <-chBlockSize:
			log.Print("set block size: ", maxSpeed)
			blockSize = getBlockSize(maxSpeed)
			break
		case <-quit:
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
