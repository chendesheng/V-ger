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
	_ "github.com/mattn/go-sqlite3"
	"github.com/nightlyone/lockfile"
	"strings"
	"time"
)

var TaskDir string

func init() {
	watchers = make([]chan *Task, 0)
	// TaskDir = path.Join(util.ReadConfig("dir"), "vger.db")
	// log.Print("Task dir:", TaskDir)

	// _, err := ioutil.ReadDir(TaskDir)
	// if os.IsNotExist(err) {
	// 	os.Mkdir(TaskDir, 0777)
	// }

	// db := openDb()
	// defer db.Close()
}

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
				ElapsedTime,
				LimitSpeed,
				Speed,
				Status,
				Est,
				Subscribe,
				Original,
				Season,
				Episode` //lastPos field move to table playing, require join table

// func SetAutoshutdown(name string, onOrOff bool) {
// 	if t, err := GetTask(name); err == nil {
// 		t.Autoshutdown = onOrOff
// 		SaveTask(t)
// 	}
// }

// func taskInfoFileName(name string) string {
// 	if !strings.HasSuffix(name, ".vger-task.txt") {
// 		name = fmt.Sprint(name, ".vger-task.txt")
// 	}
// 	return path.Join(TaskDir, name)
// }

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

	t.NameHash = hashName(t.Name)
	return t
}
func GetTask(name string) (*Task, error) {
	// println("get task:", name)
	db := openDb()
	defer db.Close()
	t, err := scanTask(db.QueryRow(fmt.Sprintf(`select %s,LastPos from task left join playing on Name=Movie where Name=?`, taskColumnes), name))
	if err != nil {
		return nil, err
	} else {
		// log.Printf("%v", t)
		return t, nil
	}
}
func ExistsEpisode(subscribeName string, season, episode int) (bool, error) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()
	var count int
	err := db.QueryRow("select count(*) from task where Subscribe=? and Season=? and Episode=?",
		subscribeName, season, episode).Scan(&count)
	return count > 0, err
}

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", TaskDir)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type taskScanner interface {
	Scan(...interface{}) error
}

func scanTask(scanner taskScanner) (*Task, error) {
	var lastPlaying sql.NullInt64

	var t Task
	var elapsedTime, est int64
	err := scanner.Scan(&t.Name,
		&t.URL,
		&t.Size,
		&t.StartTime,
		&t.DownloadedSize,
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
		log.Print(err)
		return nil, err
	}
}
func GetTasks() []*Task {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()
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

			// println(t.LastPlaying)
		}
	}

	return tasks
}

func GetDownloadingTask() (*Task, bool) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()
	t, err := scanTask(db.QueryRow(fmt.Sprintf(`select %s,LastPos from task left join playing on Name=Movie where Status='Downloading'`, taskColumnes)))
	if err != nil {
		return nil, false
	} else {
		return t, true
	}
}
func HasDownloadingOrPlaying() bool {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()
	var count int
	db.QueryRow("select count(*) from task where Statue='Downloading' or Status='Playing'").Scan(&count)

	return count > 0
}
func Exists(name string) (bool, error) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()
	var count int
	err := db.QueryRow("select count(*) from task where Name=?", name).Scan(&count)

	return count > 0, err
}

func updateTask(t *Task) error {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	_, err := db.Exec(`update task set
		URL=?,
		Size=?,
		StartTime=?,
		DownloadedSize=?,
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
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	_, err := db.Exec(fmt.Sprintf(`
			insert into task(%s) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?)
			`, taskColumnes), t.Name,
		t.URL,
		t.Size,
		t.StartTime,
		t.DownloadedSize,
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
		go writeChangeEvent(t.Name)
	} else {
		log.Print(err)
	}

	return
}
