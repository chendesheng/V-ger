package task

import (
	"encoding/base64"
	// "sync"
	// "encoding/json"
	// "errors"
	"fmt"
	// "util"
	// "io"
	// "io/ioutil"
	"log"
	// "native"
	// "net/http"
	// "os"
	// "path"
	// "strconv"
	// "regexp"
	"database/sql"
	"dbHelper"
	"strings"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNoTask = sql.ErrNoRows

type Task struct {
	URL  string
	Size int64
	Name string //identifier (a little unsafe but more readable than url)
	// seconds from 1970-1-1
	StartTime int64

	DownloadedSize   int64
	BufferedPosition int64 //only for playing task
	ElapsedTime      time.Duration

	LimitSpeed int64
	Speed      float64
	Status     string
	NameHash   string
	Est        time.Duration

	LastPlaying time.Duration

	Original  string
	Subscribe string

	Season  int
	Episode int
}

var taskColumnes string = `Name, 
				URL,
				Size,
				StartTime,
				DownloadedSize,
				BufferedPosition,
				ElapsedTime,
				LimitSpeed,
				Speed,
				Status,
				Est,
				Subscribe,
				Original,
				Season,
				Episode`

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
	t.Original = ""
	t.Subscribe = ""
	t.BufferedPosition = 0

	t.NameHash = hashName(t.Name)
	return t
}
func GetTask(name string) (*Task, error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	t, err := scanTask(db.QueryRow(fmt.Sprintf(`select %s,LastPos from task left join playing on Name=Movie where Name=?`, taskColumnes), name))
	if err != nil {
		return nil, err
	} else {
		// log.Printf("%v", t)
		return t, nil
	}
}
func GetEpisodeTask(subscribeName string, season, episode int) (string, string, string, error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	var name string
	var url string
	var status string
	err := db.QueryRow("select Name, Status, URL from task where Subscribe=? and Season=? and Episode=? and Status<>'Deleted' and Status<>'New'",
		subscribeName, season, episode).Scan(&name, &status, &url)

	return name, status, url, err
}
func ExistsEpisode(subscribeName string, season, episode int) (bool, error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	var count int
	err := db.QueryRow("select count(*) from task where Subscribe=? and Season=? and Episode=?",
		subscribeName, season, episode).Scan(&count)
	return count > 0, err
}

func scanTask(scanner dbHelper.RowScanner) (*Task, error) {
	var lastPlaying sql.NullInt64

	var t Task
	var elapsedTime, est int64
	err := scanner.Scan(&t.Name,
		&t.URL,
		&t.Size,
		&t.StartTime,
		&t.DownloadedSize,
		&t.BufferedPosition,
		&elapsedTime,
		&t.LimitSpeed,
		&t.Speed,
		&t.Status,
		&est,
		&t.Subscribe,
		&t.Original,
		&t.Season,
		&t.Episode,
		&lastPlaying)
	if err == nil {
		t.ElapsedTime = time.Duration(elapsedTime)
		t.Est = time.Duration(est)
		if lastPlaying.Valid {
			t.LastPlaying = time.Duration(lastPlaying.Int64)
		}
		return &t, nil
	} else {
		return nil, err
	}
}
func GetTasks() []*Task {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	rows, err := db.Query(fmt.Sprintf(`select %s,LastPos from task left join playing on Name=Movie`, taskColumnes))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t, err := scanTask(rows)
		if err == nil {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func GetDownloadingTask() (*Task, bool) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	t, err := scanTask(db.QueryRow(fmt.Sprintf(`select %s,LastPos from task left join playing on Name=Movie where Status='Downloading'`, taskColumnes)))
	if err != nil {
		return nil, false
	} else {
		return t, true
	}
}
func HasDownloadingOrPlaying() bool {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	var count int
	db.QueryRow("select count(*) from task where Statue='Downloading' or Status='Playing'").Scan(&count)

	return count > 0
}
func Exists(name string) (bool, error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	var count int
	err := db.QueryRow("select count(*) from task where Name=?", name).Scan(&count)

	return count > 0, err
}

func updateTask(t *Task) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec(`update task set
		URL=?,
		Size=?,
		StartTime=?,
		DownloadedSize=?,
		BufferedPosition=?,
		ElapsedTime=?,
		LimitSpeed=?,
		Speed=?,
		Status=?,
		Est=?,
		Subscribe=?,
		Original=?,
		Season=?,
		Episode=? where Name=?`,
		t.URL,
		t.Size,
		t.StartTime,
		t.DownloadedSize,
		t.BufferedPosition,
		int64(t.ElapsedTime),
		t.LimitSpeed,
		t.Speed,
		t.Status,
		int64(t.Est),
		t.Subscribe,
		t.Original,
		t.Season,
		t.Episode,
		t.Name)

	return err
}

func insertTask(t *Task) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec(fmt.Sprintf(`
			insert into task(%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?)
			`, taskColumnes), t.Name,
		t.URL,
		t.Size,
		t.StartTime,
		t.DownloadedSize,
		t.BufferedPosition,
		t.ElapsedTime,
		t.LimitSpeed,
		t.Speed,
		t.Status,
		t.Est,
		t.Subscribe,
		t.Original,
		t.Season,
		t.Episode)

	return err
}
func SaveTaskIgnoreErr(t *Task) {
	if err := SaveTask(t); err != nil {
		log.Print(err)
	}
}
func SaveTask(t *Task) (err error) {
	b, err := Exists(t.Name)
	if err != nil {
		return err
	}

	if b {
		err = updateTask(t)
	} else {
		err = insertTask(t)
	}

	if err == nil {
		go writeChangeEvent(t)
	} else {
		log.Print(err)
	}

	return
}

func writeChangeEventName(name string) {
	t, err := GetTask(name)
	if err == nil {
		writeChangeEvent(t)
	}
}

func NumOfDownloadingTasks() int {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	cnt := 0
	err := db.QueryRow("select count(*) from task where Status='Downloading'").Scan(&cnt)
	if err != nil {
		log.Print(err)
		return 0
	}
	return cnt
}

func UpdateDownloadedSize(name string, size int64) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("update task set DownloadedSize=?", size)
	if err == nil {
		go writeChangeEventName(name)
	} else {
		log.Print(err)
	}
	return err
}
