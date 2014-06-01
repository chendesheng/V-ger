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

func (gf *generateFilter) active() {
	generateBlock(gf.input, gf.output, gf.chBlockSize, gf.from, gf.to, gf.blockSize, gf.maxConnections, gf.quit)
}

func generateBlock(input chan block, output chan<- block, chBlockSize chan int64, from, size int64, blockSize int64, maxConnections int, quit <-chan bool) {
	log.Printf("generate block output: %v", output)
	if blockSize == 0 {
		//small blocksize for fast boot
		blockSize = int64(128 * 1024)
	}

	if from+blockSize > size {
		blockSize = size - from
	}

	//boot
	for i := 0; i < maxConnections; i++ {
		select {
		case output <- block{from, make([]byte, blockSize)}:
			// trace(fmt.Sprint("generate filter write boot output:", from, to))

			from += blockSize
			blockSize = blockSize
			if from+blockSize > size {
				blockSize = size - from
			}
			break
		case <-quit:
			return
		}

		if from == size {
			close(output)
			break
		}
	}

	//change to a larger blocksize after boot
	blockSize = 512 * 1024

	maxSpeed := int64(0)
	for {
		select {
		case b, ok := <-input:
			if !ok {
				return
			} else {
				if from == size {
					break
				}
				// trace(fmt.Sprint("generate filter input:", b.from, b.to))

				b.reset(from, int(blockSize))
				select {
				case output <- b:
					// trace(fmt.Sprint("generate filter write output:", b.from, b.to))

					if b.from+int64(len(b.data)) == size {
						log.Println("return generate block ", size)
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
						from += blockSize
						if from+blockSize > size {
							blockSize = size - from
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
