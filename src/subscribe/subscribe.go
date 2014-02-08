package subscribe

import (
	"database/sql"
	// "download"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	// "task"
	// "thunder"
	// "time"
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
func GetBannerImage(name string) []byte {
	db := openDb()
	defer db.Close()

	bytes := make([]byte, 0)
	err := db.QueryRow(`select BannerImage from subscribe where Name=?`, name).Scan(&bytes)
	if err != nil {
		log.Print(err)
	}

	return bytes
}
func SaveBannerImage(name string, bytes []byte) {
	db := openDb()
	defer db.Close()

	_, err := db.Exec("update subscribe set BannerImage=? where Name=?", bytes, name)
	if err != nil {
		log.Print(err)
	}
}
func updateSubscribe(s *Subscribe) error {
	db := openDb()
	defer db.Close()

	autodownload := 0
	if s.Autodownload {
		autodownload = 1
	}

	_, err := db.Exec(`update subscribe set
		Source=?,
		URL=?,
		Autodownload=?,
		Banner=? where Name=?`,
		s.Source,
		s.URL,
		autodownload,
		s.Banner,
		s.Name)

	return err
}
func insertSubscribe(s *Subscribe) error {
	db := openDb()
	defer db.Close()

	autodownload := 0
	if s.Autodownload {
		autodownload = 1
	}

	_, err := db.Exec(fmt.Sprintf(`
			insert into subscribe(%s) values(?, ?, ?, ?, ?)
			`, subscribeColumnes), s.Name,
		s.Source,
		s.URL,
		autodownload,
		s.Banner)

	return err
}
func Exists(name string) (bool, error) {
	db := openDb()
	defer db.Close()
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
