package task

import (
	"fmt"
	"time"
)

func StartNewTask(name string, url string, size int64) error {
	t := newTask(name, url, size)

	if numOfDownloadingTasks() == 0 {
		t.Status = "Downloading"
	} else {
		t.Status = "Queued"
	}

	return SaveTask(t)
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

		if numOfDownloadingTasks() == 0 {
			t.Status = "Downloading"
		} else {
			t.Status = "Queued"
		}

		return SaveTask(t)
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
	return nil
}

func ResumeNextTask() error {
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
		return SaveTask(nextTask)
	}

	return nil
}

func LimitSpeed(name string, speed int) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else {
		t.LimitSpeed = speed
		return SaveTask(t)
	}
}

func numOfDownloadingTasks() int {
	n := 0
	for _, t := range GetTasks() {
		if t.Status == "Downloading" {
			n++
		}
	}
	return n
}
