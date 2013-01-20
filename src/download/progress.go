package download

import (
	"fmt"
	"time"
)

func printProgress(progress chan int64, t *Task) {
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	part := int64(0)
	parts := [5]int64{0, 0, 0, 0, 0}
	checkTimes := [5]time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}
	cnt := 0
	percentage := float64(total) / float64(size) * 100
	speed := float64(0) // average speed of recent 5 seconds

	for length := range progress {
		total += length
		part += length

		if time.Since(checkTimes[cnt]) > time.Second || total == size {
			t.DownloadedSize = total
			elapsedTime += time.Since(checkTimes[cnt])
			t.ElapsedTime = elapsedTime
			saveTask(t)

			percentage = float64(total) / float64(size) * 100

			cnt++
			cnt = cnt % 5

			sinceLastCheck := time.Since(checkTimes[cnt])

			checkTimes[cnt] = time.Now()
			parts[cnt] = part
			part = 0

			//sum up download size of recent 5 seconds
			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			speed = float64(sum) * float64(time.Second) / float64(sinceLastCheck) / 1024
			est := time.Duration(float64((size-total))/speed) * time.Millisecond

			fmt.Printf("\r%.2f%%\t%.2f KB/s\t%s\tEst. %s     ", percentage, speed, elapsedTime, est)
		}
	}
}
