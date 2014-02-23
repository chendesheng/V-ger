package filelock

import (
	"os"
	"sync"
	"syscall"
)

type FLock struct {
	*os.File
	sync.Mutex
}

func New(filename string) (*FLock, error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.OpenFile(filename, os.O_CREATE, 0666)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &FLock{f, sync.Mutex{}}, nil
}

func (lk *FLock) Lock() error {
	err := syscall.Flock(int(lk.Fd()), syscall.LOCK_EX)
	lk.Mutex.Lock()
	return err
}

func (lk *FLock) Unlock() error {
	// f := (*os.File)(lk)
	//f.Close()
	lk.Mutex.Unlock()
	return syscall.Flock(int(lk.Fd()), syscall.LOCK_EX|syscall.LOCK_UN)
}

var DefaultLock *FLock

func Lock() {
	if DefaultLock != nil {
		DefaultLock.Lock()
	}
}

func Unlock() {
	if DefaultLock != nil {
		DefaultLock.Unlock()
	}
}
