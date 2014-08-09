package subscribe

import (
	// "database/sql"
	"sync"
	"time"
	// "download"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
	// "task"
	// "thunder"
	// "time"
	"dbHelper"
)

type Subscribe struct {
	Name         string
	Source       string
	URL          string
	Autodownload bool
	Banner       string
	Duration     time.Duration
}

var subscribeColumnes string = `Name, 
				Source,
				URL,
				Autodownload,
				Banner,
				Duration`

func scanSubscribe(scanner dbHelper.RowScanner) (*Subscribe, error) {
	var s Subscribe
	var autodownoad int
	var duration int64
	err := scanner.Scan(&s.Name,
		&s.Source,
		&s.URL,
		&autodownoad,
		&s.Banner,
		&duration)
	if err == nil {
		s.Autodownload = (autodownoad != 0)
		s.Duration = time.Duration(duration)
		return &s, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetSubscribes() []*Subscribe {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

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
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	rows, err := db.Query(fmt.Sprintf(`select %s from subscribe where Name=?`, subscribeColumnes), name)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		s, err := scanSubscribe(rows)
		if err == nil {
			log.Print("get subscribe:", s.Name)
			return s
		} else {
			log.Print(err.Error())
		}
	}

	return nil
}

func RemoveSubscribe(name string) (err error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err.Error())
		}

		return
	}()

	_, err = tx.Exec(`delete from subscribe where Name=?`, name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`update task set Status='Deleted', subscribe='' where subscribe=?`, name)
	if err != nil {
		return err
	}

	return nil
}

var bannerCache map[string][]byte
var bannerCacheLock sync.RWMutex

func GetBannerImage(name string) (bytes []byte) {
	bannerCacheLock.RLock()
	if bannerCache == nil {
		bannerCache = make(map[string][]byte)
	}

	if data, ok := bannerCache[name]; ok {
		bytes = data

		bannerCacheLock.RUnlock()
		return
	}
	bannerCacheLock.RUnlock()

	defer func() {
		if len(bytes) > 0 {
			bannerCacheLock.Lock()
			bannerCache[name] = bytes
			bannerCacheLock.Unlock()
		}
	}()
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	bytes = make([]byte, 0)
	err := db.QueryRow(`select BannerImage from subscribe where Name=?`, name).Scan(&bytes)
	if err != nil {
		log.Print(err)
	}

	return bytes
}
func SaveBannerImage(name string, bytes []byte) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("update subscribe set BannerImage=? where Name=?", bytes, name)
	if err != nil {
		log.Print(err)
	}
}
func updateSubscribe(s *Subscribe) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	autodownload := 0
	if s.Autodownload {
		autodownload = 1
	}

	_, err := db.Exec(`update subscribe set
		Source=?,
		URL=?,
		Autodownload=?,
		Banner=?,
		Duration=? where Name=?`,
		s.Source,
		s.URL,
		autodownload,
		s.Banner,
		int64(s.Duration),
		s.Name)

	return err
}
func insertSubscribe(s *Subscribe) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	autodownload := 0
	if s.Autodownload {
		autodownload = 1
	}

	_, err := db.Exec(fmt.Sprintf(`
			insert into subscribe(%s) values(?, ?, ?, ?, ?,?)
			`, subscribeColumnes), s.Name,
		s.Source,
		s.URL,
		autodownload,
		s.Banner,
		s.Duration)

	return err
}
func Exists(name string) (bool, error) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)
	var count int
	err := db.QueryRow("select count(*) from subscribe where Name=?", name).Scan(&count)

	return count > 0, err
}
func SaveSubscribe(s *Subscribe) (err error) {
	b, err := Exists(s.Name)
	if err != nil {
		return err
	}

	if b {
		err = updateSubscribe(s)
	} else {
		err = insertSubscribe(s)
	}

	if err != nil {
		log.Print(err)
	}

	return
}
func UpdateDuration(name string, duration time.Duration) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	println("update duration:", name, duration)

	_, err := db.Exec(`update subscribe set
		Duration=? where Name=?`,
		duration, name)

	return err
}
