package task

import (
	"fmt"
	"os"
	"testing"
)

func removeTask(name string) error {
	err := os.Remove(taskInfoFileName(name))
	if err != nil {
		fmt.Printf("Remove task [%s] failed: %s\n", name, err)
		return err
	}

	return nil
}

func TestTaskDir(t *testing.T) {
	tk := &Task{}
	tk.Name = "hello"
	SaveTask(tk)
	tk2, _ := GetTask(tk.Name)

	if fmt.Sprintf("%v", tk) != fmt.Sprintf("%v", tk2) {
		t.Errorf("not equal tk is %v, tk2 is %v", tk, tk2)
	}

	removeTask(tk.Name)

	_, err := GetTask(tk.Name)
	if err == nil {
		t.Error("should error")
	}
}

func TestGetTasks(t *testing.T) {
	tk := &Task{}
	tk.Name = "hello"
	SaveTask(tk)

	tk = &Task{}
	tk.Name = "hello1"
	SaveTask(tk)

	ts := GetTasks()
	if len(ts) != 2 {
		t.Errorf("length must equals 2 now equals %d", len(ts))
	}

	removeTask(tk.Name)
	removeTask("hello")
}

func TestWatchTasksChange(t *testing.T) {
	ch := make(chan *Task)
	WatchChange(ch)
	WatchChange(ch)

	go func() {
		tk := <-ch
		if tk.Name != "hello" {
			t.Errorf("Except hello now %s", tk.Name)
		}
		tk = <-ch
		if tk.Name != "hello1" {
			t.Errorf("Except hello1 now %s", tk.Name)
		}
	}()

	tk := &Task{}
	tk.Name = "hello"
	SaveTask(tk)

	tk = &Task{}
	tk.Name = "hello1"
	SaveTask(tk)

	GetTasks()

	RemoveWatch(ch)
	removeTask("hello1")
	removeTask("hello")
}

//TODO
func TestMultiWatcher(t *testing.T) {

}

func TestGetDownloadingTask(t *testing.T) {
	tk := &Task{}
	tk.Name = "hello3"
	tk.Status = "Downloading"
	SaveTask(tk)

	tk2, _ := GetDownloadingTask()

	if fmt.Sprintf("%v", tk) != fmt.Sprintf("%v", tk2) {
		t.Errorf("Except %v but %v", tk, tk2)
	}

	tk2.Status = "Stopped"
	SaveTask(tk2)

	tk3, ok := GetDownloadingTask()

	if ok {
		t.Errorf("Except <nil> but %v", tk3)
	}

	removeTask(tk.Name)
}

func TestCleanup(t *testing.T) {
	os.Remove(TaskDir)
}
