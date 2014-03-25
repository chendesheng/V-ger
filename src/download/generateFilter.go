package download

import (
	"fmt"
	"log"
	"time"
	"util"
)

type generateFilter struct {
	basicFilter
	from        int64
	to          int64
	blockSize   int64
	chBlockSize chan int64
}

func (gf *generateFilter) active() {
	generateBlock(gf.input, gf.output, gf.chBlockSize, gf.from, gf.to, gf.blockSize, gf.quit)
}

type blockPool struct {
	data []*block
	max  int
}

func (p *blockPool) get(from, to int64) *block {
	if len(p.data) == 0 {
		println("new block:", from, to)
		return &block{from, to, make([]byte, 0, to-from)}
	} else {
		lastIndex := len(p.data) - 1
		b := p.data[lastIndex]
		p.data = p.data[:lastIndex]
		println("get from pool:", from, to, cap(b.data))

		b.from, b.to, b.data = from, to, b.data[:0]
		return b
	}
	// return &block{from, to, make([]byte, 0, to-from+bytes.MinRead)}
}
func (p *blockPool) put(b *block) {
	if len(p.data) < p.max {
		p.data = append(p.data, b)
	}
}

func generateBlock(input chan *block, output chan<- *block, chBlockSize chan int64, from, size int64, blockSize int64, quit <-chan bool) {
	log.Printf("generate block output: %v", output)
	if blockSize == 0 {
		blockSize = int64(400 * 1024)
	}

	to := from + blockSize
	if to > size {
		to = size
	}

	//small blocksize after start,
	//change to a larger blocksize after 15 seconds
	changeBlockSize := time.NewTimer(time.Second * 15)
	startCnt := util.ReadIntConfig("max-connection")

	p := blockPool{}
	p.data = make([]*block, 0, startCnt)
	p.max = startCnt

	log.Printf("output %v", output)
	maxSpeed := int64(0)
	for {
		if startCnt < 0 {
			select {
			case b, ok := <-input:
				if !ok {
					return
				} else {
					p.put(b)
				}
			case <-quit:
				return
			}
		} else {
			startCnt--
		}

		// b := time.Now()
		b := p.get(from, to)

		select {
		case maxSpeed = <-chBlockSize:
			// maxSpeed = 0
			log.Print("set block size: ", maxSpeed)
			if maxSpeed > 0 {
				blockSize = maxSpeed * 1024
			} else {
				blockSize = int64(32 * 1024)
				changeBlockSize.Reset(time.Second * 15)
			}
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
		case <-changeBlockSize.C:
			if maxSpeed == 0 {
				blockSize = 400 * 1024
			}
			changeBlockSize.Stop()
		case <-quit:
			// close(output)
			fmt.Println("quit generate block")
			return
		}
	}
}
