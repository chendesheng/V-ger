package download

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Task struct {
	URL       string
	Size      int64
	Name      string //identifier (a little unsafe but more readable than url)
	Path      string
	StartDate string

	DownloadedSize int64
	ElapsedTime    time.Duration
}

func taskInfoFileName(taskName string) string {
	return fmt.Sprintf("tasks%c%s.vger-task.txt", os.PathSeparator, taskName)
}
func saveTask(t Task) {
	writeJson(taskInfoFileName(t.Name), t)
}
func removeTask(name string) {
	err := os.Remove(taskInfoFileName(name))
	if err != nil {
		fmt.Printf("Remove task [%s] failed: %s\n", name, err)
	}
}
func getOrNewTask(url string, name string) (Task, bool) {
	for _, t := range getTasks() {
		if name == t.Name {
			return t, false
		}
	}

	t := Task{URL: url, Name: name}
	return t, true
}

//one second cache for task list
var taskCache []Task //TODO: need lock

func getTasks() []Task {
	if taskCache != nil {
		return taskCache
	}

	fileInfoes, err := ioutil.ReadDir("tasks")
	if err != nil {
		log.Fatal(err)
	}

	tasks := make([]Task, len(fileInfoes))
	for _, f := range fileInfoes {
		name := f.Name()
		if f.IsDir() || !strings.HasSuffix(name, ".vger-task.txt") {
			continue
		}

		t := Task{}
		readJson(name, t)
		tasks = append(tasks, t)
	}

	taskCache = tasks
	chanTimeout := time.Tick(time.Second)
	go func() {
		<-chanTimeout
		taskCache = nil
	}()
	log.Println(tasks)
	return tasks
}
