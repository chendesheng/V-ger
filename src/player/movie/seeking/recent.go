package seeking

func recentPipe(quit chan struct{}) (in chan *seekArg, out chan *seekArg) {
	in = make(chan *seekArg)
	out = make(chan *seekArg)

	go func() {
		var recentValue *seekArg
		var outnil chan *seekArg
		for {
			select {
			case t, ok := <-in:
				if !ok {
					return
				}
				outnil = out
				recentValue = t
			case outnil <- recentValue:
				outnil = nil
			case <-quit:
				return
			}
		}
	}()

	return
}
