package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"logger"
	"os"
	"path"
	"strings"
	"time"
	"util"
)

var TaskDir string

// func init() {
// 	watchers = make([]chan *Task, 0)
// 	TaskDir = path.Join(util.ReadConfig("dir"), "vger-tasks")
// 	log.Print("Task dir:", TaskDir)

// 	_, err := ioutil.ReadDir(TaskDir)
// 	if os.IsNotExist(err) {
// 		os.Mkdir(TaskDir, 0777)
// 	}
// }

type Task struct {
	URL  string
	Size int64
	Name string //identifier (a little unsafe but more readable than url)
	// seconds from 1970-1-1
	StartTime int64

	DownloadedSize int64
	ElapsedTime    time.Duration

	LimitSpeed int64
	Speed      float64
	Status     string
	NameHash   string
	Est        time.Duration

	// Autoshutdown bool

	Subs        []string
	LastPlaying time.Duration
}

// func SetAutoshutdown(name string, onOrOff bool) {
// 	if t, err := GetTask(name); err == nil {
//  		t.Autoshutdown = onOrOff
// 		SaveTask(t)
// 	}
// }

func taskInfoFileName(name string) string {
	if !strings.HasSuffix(name, ".vger-task.txt") {
		name = fmt.Sprint(name, ".vger-task.txt")
	}
	return path.Join(TaskDir, name)
}

func hashName(name string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(name)), "=")
}
func newTask(name string, url string, size int64) *Task {
	t := new(Task)
	t.URL = url
	t.Name = name
	t.Size = size
	t.StartTime = time.Now().Unix()
	t.DownloadedSize = 0
	t.ElapsedTime = 0

	t.LimitSpeed = 0
	t.Speed = 0
	t.Status = "New"

	t.NameHash = hashName(t.Name)
	return t
}
func GetTask(name string) (*Task, error) {
	t := new(Task)
	err := util.ReadJson(taskInfoFileName(name), t)
	if err != nil {
		log.Printf("Get task error:%s. Task name:%s.", err.Error(), name)
		return nil, err
	}

	return t, nil
}

func GetTasks() []*Task {
	fileInfoes, err := ioutil.ReadDir(TaskDir)
	if err != nil {
		log.Print(err)
		return make([]*Task, 0)
	}

	tasks := make([]*Task, 0, len(fileInfoes))
	for _, f := range fileInfoes {
		name := f.Name()

		if strings.HasPrefix(name, ".") || f.IsDir() || !strings.HasSuffix(name, ".vger-task.txt") { //exculding hidden files
			continue
		}

		if t, err := GetTask(name); err == nil {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func GetDownloadingTask() (*Task, bool) {
	for _, t := range GetTasks() {
		if t.Status == "Downloading" {
			return t, true
		}
	}

	return nil, false
}
func HasDownloadingOrPlaying() bool {
	for _, t := range GetTasks() {
		if t.Status == "Downloading" || t.Status == "Playing" {
			log.Printf("has downloading or playing %v", t)
			return true
		}
	}

	return false
}
func SaveTask(t *Task) (err error) {
	err = util.WriteJson(taskInfoFileName(t.Name), t)
	if err == nil {
		// go writeChangeEvent(t.Name)
	}

	return
}

func main() {
	logger.InitLog("a.log")

	os.Remove("./vger.db")

	TaskDir = "/Volumes/Data/Downloads/Video/vger-tasks"

	db, err := sql.Open("sqlite3", "./vger.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	create table task (
		Name NVARCHAR(2048) not null primary key default(''), 
		Size INTEGER not null default(0),
		URL NVARCHAR(2048) not null default(''),
		StartTime INTEGER not null default(0),
		DownloadedSize INTEGER not null default(0),
		ElapsedTime INTEGER not null default(0),
		LimitSpeed INTEGER not null default(0),
		Speed DOUBLE not null default(0),
		Status NVARCHAR(128) not null default(0),
		Est INTEGER not null default(0),
		LastPlaying INTEGER not null default(0)
--		Sub1Path text not null default(''),
--		Sub1Offset INTEGER not null default(0),
--		Sub2Path text not null default(''),
--		Sub2Offset INTEGER not null default(0)
		);
	delete from task;

	create table 
	`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}

	for _, t := range GetTasks() {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare(`
			insert into task(
				Name, 
				URL,
				Size,
				StartTime,
				DownloadedSize,
				ElapsedTime,
				LimitSpeed,
				Speed,
				Status,
				Est,
				LastPlaying
				) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(t.Name,
			t.URL,
			t.Size,
			t.StartTime,
			t.DownloadedSize,
			t.ElapsedTime,
			t.LimitSpeed,
			t.Speed,
			t.Status,
			t.Est,
			t.LastPlaying)
		if err != nil {
			log.Fatal(err)
		}

		stmt.Close()
		tx.Commit()
		if err != nil {
			log.Printf("%q: %s\n", err, sql)
		}
	}
}
