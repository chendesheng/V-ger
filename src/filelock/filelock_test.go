package filelock

import (
	"fmt"
	"testing"
	"time"
)

func TestLockFileSync(t *testing.T) {
	filename := "/tmp/lock.me.now.lck"

	lock, err := New(filename)
	if err != nil {
		fmt.Println("Cannot init lock. reason: ", err)
		panic(err)
	}

	go func() {
		lock1, err := New(filename)
		if err != nil {
			println(err)
		}
		err = lock1.Lock()
		println(0)
		time.Sleep(time.Second)
		println(0)

		lock1.Unlock()
	}()

	err = lock.Lock()
	println(1)
	time.Sleep(1 * time.Second)
	println(1)
	lock.Unlock()

	println("after unlock")
	time.Sleep(3 * time.Second)
}

func TestLockFileSync2(t *testing.T) {
	filename := "/tmp/lock.me.now.lck"

	lock, err := New(filename)
	if err != nil {
		fmt.Println("Cannot init lock. reason: ", err)
		panic(err)
	}

	go func() {
		// lock1, err := New(filename)
		if err != nil {
			println(err)
		}
		err = lock.Lock()
		println(0)
		time.Sleep(time.Second)
		println(0)

		lock.Unlock()
	}()

	err = lock.Lock()
	println(1)
	time.Sleep(1 * time.Second)
	println(1)
	lock.Unlock()

	println("after unlock")
	time.Sleep(3 * time.Second)
}
func TestLockFileSync3(t *testing.T) {
	filename := "/tmp/lock.me.now.lck"

	lock, err := New(filename)
	if err != nil {
		fmt.Println("Cannot init lock. reason: ", err)
		panic(err)
	}

	go func() {
		lock1, err := New(filename)
		if err != nil {
			println(err)
		}
		err = lock1.Lock()
		println(0)
		time.Sleep(time.Second)
		println(0)

		lock1.Unlock()
	}()

	err = lock.Lock()
	println(1)
	time.Sleep(200 * time.Millisecond)
	println(1)
	lock.Unlock()
	time.Sleep(200 * time.Millisecond)

	err = lock.Lock()
	println(2)
	time.Sleep(200 * time.Millisecond)
	println(2)
	lock.Unlock()

	println("after unlock")
	time.Sleep(3 * time.Second)
}

func TestLeak(t *testing.T) {
	filename := "/tmp/lock.leak.lck"
	lk, _ := New(filename)
	lk.Lock()
}
