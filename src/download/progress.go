package download

import (
	"fmt"
	"log"
	"time"
)

func handleProgress(progress chan int64, t *Task) {
	log.Printf("start handle progress: %v\n", *t)
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	timer := time.NewTicker(time.Second * 2)

	speed := float64(0)
	partsCount := 15
	parts := make([]int64, partsCount)
	checkTimes := make([]time.Time, partsCount)
	for i := 0; i < partsCount; i++ {
		parts[i] = 0
		checkTimes[i] = time.Now()
	}
	part := int64(0)
	cnt := 0
	est := time.Duration(0)
	lastCheck := time.Now()

	for {
		select {
		case length, ok := <-progress:
			if !ok {
				saveProgress(t.Name, speed, total, elapsedTime, 0)
				return
			}
			// fmt.Println("progress ", total)
			total += length
			part += length

			if time.Since(checkTimes[cnt]) > time.Second || total == size {
				cnt++
				cnt = cnt % partsCount

				lastCheck = checkTimes[cnt]
				checkTimes[cnt] = time.Now()
				parts[cnt] = part
				part = 0
			}
		case <-timer.C:
			elapsedTime += time.Second * 2

			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			speed = float64(sum) * float64(time.Second) / float64(time.Since(lastCheck)) / 1024

			percentage, est := calcProgress(total, size, speed)

			saveProgress(t.Name, speed, total, elapsedTime, est)

			printProgress(percentage, speed, elapsedTime, est)
			if total == size {
				fmt.Println("progress return")
				return
			}
		}
	}
	printProgress(100, speed, elapsedTime, est)
}
func calcProgress(total, size int64, speed float64) (percentage float64, est time.Duration) {
	percentage = float64(total) / float64(size) * 100
	if speed == 0 {
		est = 0
	} else {
		est = time.Duration(float64((size-total))/speed) * time.Millisecond
	}
	return
}
func saveProgress(name string, speed float64, total int64, elapsedTime time.Duration, est time.Duration) {
	if t, ok := GetTask(name); ok {
		t.DownloadedSize = total
		t.ElapsedTime = elapsedTime
		t.Speed = speed
		t.Est = est
		saveTask(t)
	}
}
func printProgress(percentage float64, speed float64, elapsedTime time.Duration, est time.Duration) {
	// fmt.Printf("%.2f%%    %.2f KB/s    %s    Est. %s     \n", percentage, speed, elapsedTime/time.Second*time.Second, est/time.Second*time.Second)
}
