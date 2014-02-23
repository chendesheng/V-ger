package task

import (
	// "fmt"
	// "log"
	"sync"
	// "time"
)

var watchers []chan *Task
var watchersLock sync.Mutex = sync.Mutex{}

func WatchChange(ch chan *Task) {
	if ch == nil {
		panic("ch cannot be nil")
	}

	watchersLock.Lock()
	defer watchersLock.Unlock()

	for _, w := range watchers {
		if w == ch {
			return
		}
	}

	watchers = append(watchers, ch)
	// chTaskChange = ch
}

func RemoveWatch(ch chan *Task) {
	watchersLock.Lock()
	defer watchersLock.Unlock()

	for i, w := range watchers {
		if w == ch {
			if i == len(watchers)-1 {
				watchers = watchers[:i]
			} else {
				watchers = append(watchers[:i], watchers[i+1:]...)
			}

			// log.Println("remove watch: ", w)
			break
		}
	}
}

func writeChangeEvent(t *Task) {
	watchersLock.Lock()
	defer watchersLock.Unlock()

	for _, w := range watchers {
		go func(w chan *Task) {
			w <- t
		}(w)
	}
}
