package subscribe

import (
	"database/sql"
	"download"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"task"
	"thunder"
	"time"
)

type Subscribe struct {
	Name         string
	Source       string
	URL          string
	Autodownload bool
	Banner       string
}

var subscribeColumnes string = `Name, 
				Source,
				URL,
				Autodownload,
				Banner`
var DbPath string

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type myScanner interface {
	Scan(...interface{}) error
}

func scanSubscribe(scanner myScanner) (*Subscribe, error) {
	var s Subscribe
	var autodownoad int
	err := scanner.Scan(&s.Name,
		&s.Source,
		&s.URL,
		&autodownoad,
		&s.Banner)
	if err == nil {
		s.Autodownload = (autodownoad != 0)
		return &s, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetSubscribes() []*Subscribe {
	db := openDb()
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf(`select %s from subscribe`, subscribeColumnes))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	subscribes := make([]*Subscribe, 0)
	for rows.Next() {
		s, err := scanSubscribe(rows)
		if err == nil {
			subscribes = append(subscribes, s)
		}
	}

	return subscribes
}
func GetSubscribe(name string) *Subscribe {
	db := openDb()
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf(`select %s from subscribe where Name=?`, subscribeColumnes), name)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		s, err := scanSubscribe(rows)
		if err == nil {
			return s
		}
	}

	return nil
}
func SaveSubscribe(s *Subscribe) (err error) {
	db := openDb()
	defer db.Close()

	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	var count int
	err = tx.QueryRow("select count(*) from subscribe where Name=?", s.Name).Scan(&count)
	if err != nil {
		return err
	}

	var autodownload int
	autodownload = 0
	if s.Autodownload {
		autodownload = 1
	}

	if count > 0 {
		// println("update")

		_, err = tx.Exec(`update subscribe set
		Source=?,
		URL=?,
		Autodownload=?,
		Banner=? where Name=?`,
			s.Source,
			s.URL,
			autodownload,
			s.Banner,
			s.Name)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err = tx.Exec(fmt.Sprintf(`
			insert into subscribe(%s) values(?, ?, ?, ?, ?)
			`, subscribeColumnes), s.Name,
			s.Source,
			s.URL,
			autodownload,
			s.Banner)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return
}

func UpdateAll() {
	subscribes := GetSubscribes()
	for _, s := range subscribes {
		_, tasks, err := Parse(s.URL)
		if err != nil {
			log.Print(err)
		} else {
			for _, t := range tasks {
				if b, err := task.Exists(t.Name); err == nil && !b {
					log.Printf("subscribe new task: %v", t)

					if t.Season < 0 {
						task.SaveTask(t)
						continue
					}

					files, err := thunder.NewTask(t.Original, "")
					if err != nil {
						log.Print(err)
					}
					fmt.Printf("%v\n", files)
					if err == nil && len(files) == 1 && files[0].Percent == 100 {
						t.URL = files[0].DownloadURL
						_, _, size, err := download.GetDownloadInfo(t.URL)
						if err != nil {
							log.Print(err)
						} else {
							t.Size = size
							task.SaveTask(t)
							task.StartNewTask2(t)

						}
					}
				}
			}
		}
	}
}

func Monitor() {
	time.Sleep(3 * time.Second)

	for {
		UpdateAll()

		time.Sleep(30 * time.Second)
	}
}
