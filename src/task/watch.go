package task

import (
	"log"
	"sync"
	"time"
)

var watchers []chan *Task
var watcherLock sync.Mutex = sync.Mutex{}

func WatchChange(ch chan *Task) {
	if ch == nil {
		panic("ch cannot be nil")
	}

	watcherLock.Lock()
	defer watcherLock.Unlock()

	for _, w := range watchers {
		if w == ch {
			return
		}
	}

	watchers = append(watchers, ch)
	// chTaskChange = ch
}

func RemoveWatch(ch chan *Task) {
	watcherLock.Lock()
	defer watcherLock.Unlock()

	for i, w := range watchers {
		if w == ch {
			if i == len(watchers)-1 {
				watchers = watchers[:i]
			} else {
				watchers = append(watchers[:i], watchers[i+1:]...)
			}

			log.Println("remove watch: ", w)
			break
		}
	}
}

//call this function after modify task file directly, like trash task.
// func UpdateFiles() {
// 	writeChangeEvent()
// }

func writeChangeEvent(name string) {
	// tks := GetTasks()

	t, err := GetTask(name)
	if err != nil {
		log.Println(err)
		return
	}

	watcherLock.Lock()
	copyWatchers := make([]chan *Task, len(watchers))
	copy(copyWatchers, watchers)
	watcherLock.Unlock()

	for _, w := range copyWatchers {
		select {
		case w <- t:
			break
		case <-time.After(time.Second):
			log.Printf("writeChangeEvent timeout: %v\n", w)
			break
		}
	}
}
