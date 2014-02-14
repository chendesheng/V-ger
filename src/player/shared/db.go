package shared

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"

	"github.com/nightlyone/lockfile"
)

var DbFile = ""

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", DbFile)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type rowScanner interface {
	Scan(...interface{}) error
}

func scanSub(scanner rowScanner) (*Sub, error) {
	var sub Sub

	var offset int64
	err := scanner.Scan(&sub.Movie, &sub.Name, &offset, &sub.Content, &sub.Type, &sub.Lang1, &sub.Lang2)
	if err == nil {
		sub.Offset = time.Duration(offset)
		return &sub, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetSubtitles(movie string) []*Sub {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	println("get local subtitles:", movie)

	sql := `select Movie, Name, Offset, Content, Type, Lang1, Lang2 from subtitle where Movie=?`
	rows, err := db.Query(sql, movie)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	subs := make([]*Sub, 0)
	for rows.Next() {
		sub, err := scanSub(rows)
		if err == nil {
			subs = append(subs, sub)
		}
	}

	return subs
}
func GetSubtitle(name string) *Sub {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	sql := `select Movie, Name, Offset, Content, Type, Lang1, Lang2 from subtitle where Name=?`
	rows, err := db.Query(sql, name)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		sub, err := scanSub(rows)
		if err == nil {
			return sub
		}
	}

	return nil
}
func InsertSubtitle(sub *Sub) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	var count int
	err := db.QueryRow("select count(*) from subtitle where Name=?", sub.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		sql := "insert into subtitle(Movie, Name, Offset, Content, Type, Lang1, Lang2) values (?,?,?,?,?,?,?)"
		_, err := db.Exec(sql, sub.Movie, sub.Name, sub.Offset, sub.Content, sub.Type, sub.Lang1, sub.Lang2)
		if err != nil {
			log.Print(err)
		}
	}
}

func UpdateSubtitleOffset(name string, offset time.Duration) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	_, err := db.Exec("update subtitle set Offset=? where Name=?", int64(offset), name)
	if err != nil {
		log.Print(err)
	}
}

func UpdateSubtitleLanguage(name string, lang1, lang2 string) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	_, err := db.Exec("update subtitle set Lang1=?,Lang2=? where Name=?", lang1, lang2, name)
	if err != nil {
		log.Print(err)
	}
}

func scanPlaying(scanner rowScanner) (*Playing, error) {
	var p Playing
	var lastPos, duration int64
	err := scanner.Scan(&p.Movie, &lastPos, &p.SoundStream, &p.Sub1, &p.Sub2, &duration)
	if err == nil {
		p.LastPos = time.Duration(lastPos)
		p.Duration = time.Duration(duration)
		return &p, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetPlaying(movie string) *Playing {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	db := openDb()
	defer db.Close()

	sql := `select Movie, LastPos, SoundStream, Sub1, Sub2, Duration from Playing where Movie=?`
	rows, err := db.Query(sql, movie)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		p, err := scanPlaying(rows)

		if err != nil {
			log.Print(err)
			return nil
		} else {
			return p
		}
	} else {
		return nil
	}
}

func CreateOrGetPlaying(movie string) *Playing {
	db := openDb()
	defer db.Close()

	var count int
	err := db.QueryRow("select count(*) from playing where Movie=?", movie).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		SavePlaying(&Playing{movie, 0, -1, "", "", 0})
	}

	return GetPlaying(movie)
}

func SavePlaying(p *Playing) {
	if len(lockfile.DefaultLock) > 0 {
		lockfile.DefaultLock.Lock()
		defer lockfile.DefaultLock.Unlock()
	}

	log.Printf("Save playing: %v", *p)

	db := openDb()
	defer db.Close()

	var count int
	err := db.QueryRow("select count(*) from playing where Movie=?", p.Movie).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		sql := "insert into playing(Movie, LastPos, SoundStream, Sub1, Sub2, Duration) values (?,?,?,?,?,?)"
		_, err := db.Exec(sql, p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2, p.Duration)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := db.Exec("update playing set Movie=?,LastPos=?,SoundStream=?,Sub1=?,Sub2=?,Duration=? where Movie=?",
			p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2, p.Duration, p.Movie)
		if err != nil {
			log.Print(err)
		}
	}
}
