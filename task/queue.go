package task

import (
	"fmt"
	"time"
	"vger/util"
)

func NewTask(name string, url string, size int64, status string) error {
	t := newTask(name, url, size)
	t.Status = status
	return SaveTask(t)
}

func StartNewTask(name string, url string, size int64) error {
	t := newTask(name, url, size)

	startOrQueueTask(t)

	return SaveTask(t)
}
func StartNewTask2(t *Task) {
	println("start new task2")
	startOrQueueTask(t)

	SaveTask(t)
}

func ResumeTask(name string) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else if t.Status == "Finished" {
		return fmt.Errorf("The task is already finished.")
	} else if t.Status == "Downloading" {
		return nil
	} else {
		startOrQueueTask(t)
		t.Speed = 0
		return SaveTask(t)
	}
}

func startOrQueueTask(t *Task) bool {
	if NumOfDownloadingTasks() < util.ReadIntConfig("simultaneous-downloads") {
		t.Status = "Downloading"
		return true
	} else {
		t.Status = "Queued"
		return false
	}
}

func DeleteTask(name string) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else {
		if len(t.Subscribe) > 0 {
			t.Status = "New"
			t.LastPlaying = 0
			t.DownloadedSize = 0
			t.Speed = 0
		} else {
			t.Status = "Deleted"
		}

		SaveTask(t)
	}
	return nil
}

func StopTask(name string) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else if t.Status == "Finished" {
		return fmt.Errorf("The task is already finished.")
	} else {
		t.Status = "Stopped"
		return SaveTask(t)
	}
}

func ResumeNextTask() (error, bool) {
	tasks := GetTasks()

	var nextTask *Task
	startTime := time.Now().Unix()
	for _, t := range tasks {
		if t.Status == "Queued" && t.StartTime < startTime {
			startTime = t.StartTime
			nextTask = t
		}
	}
	if nextTask != nil {
		nextTask.Status = "Downloading"
		return SaveTask(nextTask), true
	} else {
		return nil, false
	}
}

func LimitSpeed(name string, speed int64) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else {
		t.LimitSpeed = speed
		return SaveTask(t)
	}
}

func QueueDownloadingTask() error {
	tasks := GetTasks()

	var latestDownloadTask *Task
	startTime := (time.Time{}).Unix()
	for _, t := range tasks {
		if t.Status == "Downloading" && t.StartTime > startTime {
			startTime = t.StartTime
			latestDownloadTask = t
		}
	}
	if latestDownloadTask != nil {
		latestDownloadTask.Status = "Queued"

		return SaveTask(latestDownloadTask)
	}

	return nil
}
