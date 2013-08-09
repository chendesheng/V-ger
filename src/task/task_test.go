package task

import (
	"fmt"
	"os"
	"testing"
)

func testSaveTask(t *testing.T) {
	tk := &Task{}
	tk.Name = "hello"
	SaveTask(tk)
	tk2, _ := GetTask(tk.Name)

	if fmt.Sprintf("%v", tk) != fmt.Sprintf("%v", tk2) {
		t.Errorf("not equal tk is %v, tk2 is %v", tk, tk2)
	}

	RemoveTask(tk.Name)

	_, err := GetTask(tk.Name)
	if err == nil {
		t.Error("should error")
	}
}

func TestTaskDir(t *testing.T) {
	taskDir = ""
	testSaveTask(t)

	taskDir = "abc"
	testSaveTask(t)

	os.Remove(taskDir)
}

func TestGetTasks(t *testing.T) {
	taskDir = "abc"

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

	RemoveTask(tk.Name)
	RemoveTask("hello")
	os.Remove(taskDir)
}

func TestWatchTasksChange(t *testing.T) {
	taskDir = "abc"

	ch := make(chan []*Task)
	WatchChange(ch)
	WatchChange(ch)

	go func() {
		tks := <-ch
		if len(tks) != 1 {
			t.Errorf("length of tasks must be 1 now %d", len(tks))
		}
		if tks[0].Name != "hello" {
			t.Errorf("Except hello now %s", tks[0].Name)
		}
		tks = <-ch
		if len(tks) != 2 {
			t.Errorf("length of tasks must be 2 now %d", len(tks))
		}
		tks = <-ch
		if len(tks) != 1 {
			t.Errorf("length of tasks must be 1 now %d", len(tks))
		}
		if tks[0].Name != "hello1" {
			t.Errorf("Except hello1 now %s", tks[0].Name)
		}
		tks = <-ch
		if len(tks) != 0 {
			t.Errorf("Except 0 but %d", len(tks))
		}
	}()

	tk := &Task{}
	tk.Name = "hello"
	SaveTask(tk)

	tk = &Task{}
	tk.Name = "hello1"
	SaveTask(tk)

	GetTasks()

	RemoveTask("hello")
	RemoveTask("hello1")

	os.Remove(taskDir)
}

//TODO
func TestMultiWatcher(t *testing.T) {

}

func TestGetDownloadingTask(t *testing.T) {
	taskDir = "abc"
	tk := &Task{}
	tk.Name = "hello"
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

	RemoveTask(tk.Name)
	os.Remove(taskDir)
}
