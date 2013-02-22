package download

import (
	"fmt"
	"log"
	"time"
)

func handleProgress(progress chan int64, t *Task) {
	log.Printf("start handle progress: %v\n", *t)
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	timer := time.After(time.Second * 2)

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

	for total < size {
		select {
		case length := <-progress:
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

		case <-timer:
			timer = time.After(time.Second * 2)

			elapsedTime += time.Second * 2
			saveProgress(t, total, elapsedTime)

			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			speed = float64(sum) * float64(time.Second) / float64(time.Since(lastCheck)) / 1024

			percentage, est := calcProgress(total, size, speed)
			printProgress(percentage, speed, elapsedTime, est)
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
func saveProgress(t *Task, total int64, elapsedTime time.Duration) {
	t.DownloadedSize = total
	t.ElapsedTime = elapsedTime
	saveTask(t)
}
func printProgress(percentage float64, speed float64, elapsedTime time.Duration, est time.Duration) {
	fmt.Printf("\r%.2f%%    %.2f KB/s    %s    Est. %s     ", percentage, speed, elapsedTime/time.Second*time.Second, est/time.Second*time.Second)
}
