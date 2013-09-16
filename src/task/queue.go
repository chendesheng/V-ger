package task

import (
	"fmt"
	"time"
	"util"
)

func StartNewTask(name string, url string, size int64) error {
	t := newTask(name, url, size)

	if numOfDownloadingTasks() < 2 {
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
		if numOfDownloadingTasks() < util.ReadIntConfig("simultaneous-downloads") {
			t.Status = "Downloading"
		} else {
			t.Status = "Queued"
		}
		t.Speed = 0
		return SaveTask(t)
	}

	return nil
}

func DeleteTask(name string) error {
	t, err := GetTask(name)
	if err != nil {
		return err
	} else {
		t.Status = "Deleted"
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
	return nil
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

func numOfDownloadingTasks() int {
	n := 0
	for _, t := range GetTasks() {
		if t.Status == "Downloading" {
			n++
		}
	}
	fmt.Println("num of downloading tasks ", n)
	return n
}
