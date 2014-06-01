package download

import "fmt"

type simpleWriteFilter struct {
	basicFilter
	w WriterAtQuit
}

func (swf *simpleWriteFilter) active() {
	defer swf.closeOutput()

	for {
		select {
		case b, ok := <-swf.input:
			if !ok {
				fmt.Println("close simple write output")
				return
			}

			// trace(fmt.Sprint("simple write filter input:", b.from, b.to))

			swf.w.WriteAtQuit(b.data, b.from, swf.quit)

			swf.writeOutput(b)
			// trace(fmt.Sprint("simple write filter output:", b.from, b.to))

			break
		case <-swf.quit:
			fmt.Println("simple write output quit")
			return
		}
	}

	fmt.Println("simpleWriteOutput end")
}
