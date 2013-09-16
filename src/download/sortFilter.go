package download

import (
	"fmt"
)

type sortFilter struct {
	basicFilter
	from int64
}

func (sf *sortFilter) active() {
	sortOutput(sf.input, sf.output, sf.quit, sf.from)
}

func sortOutput(input <-chan *block, output chan<- *block, quit <-chan bool, from int64) {
	dbmap := make(map[int64]*block)
	nextOutputFrom := from
	for {
		select {
		case db, ok := <-input:
			if db != nil {
				dbmap[db.from] = db
			}

			// log.Println(len(dbmap))
			for {
				if d, exist := dbmap[nextOutputFrom]; exist {
					// fmt.Printf("sort output %d-%d\n", d.from, d.to)
					select {
					case output <- d:
						nextOutputFrom = d.to
						delete(dbmap, d.from)
						break
					case <-quit:
						return
					}
				} else {
					break
				}
			}

			if !ok {
				close(output)
				return
			}
		case <-quit:
			fmt.Println("sort output quit")
			return
		}
	}
}
