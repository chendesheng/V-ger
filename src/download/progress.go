package download

import (
	"fmt"
	"time"
)

func printProgress(progress chan int64, t *Task) {
	size, total, elapsedTime := t.Size, t.DownloadedSize, t.ElapsedTime

	partsCount := 10
	cnt := 0
	part := int64(0)
	parts := make([]int64, partsCount)
	checkTimes := make([]time.Time, partsCount)
	for i := 0; i < partsCount; i++ {
		parts[i] = 0
		checkTimes[i] = time.Now()
	}

	for length := range progress {
		total += length
		part += length

		if time.Since(checkTimes[cnt]) > time.Second || total == size {
			t.DownloadedSize = total
			elapsedTime += time.Since(checkTimes[cnt])
			t.ElapsedTime = elapsedTime
			saveTask(t)

			cnt++
			cnt = cnt % partsCount

			sinceLastCheck := time.Since(checkTimes[cnt])

			checkTimes[cnt] = time.Now()
			parts[cnt] = part
			part = 0

			//sum up download size of recent 5 seconds
			sum := int64(0)
			for _, p := range parts {
				sum += p
			}
			percentage := float64(total) / float64(size) * 100
			speed := float64(sum) * float64(time.Second) / float64(sinceLastCheck) / 1024
			est := time.Duration(float64((size-total))/speed) * time.Millisecond

			fmt.Printf("\r%.2f%%    %.2f KB/s    %s    Est. %s     ", percentage, speed, elapsedTime/time.Second*time.Second, est/time.Second*time.Second)
		}
	}
}
