package shared

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
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
	err := scanner.Scan(&sub.Movie, &sub.Name, &offset, &sub.Content, &sub.Type)
	if err == nil {
		sub.Offset = time.Duration(offset)
		return &sub, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetSubtitles(movie string) []*Sub {
	db := openDb()
	defer db.Close()

	sql := `select Movie, Name, Offset, Content, Type from subtitle where Movie=?`
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
	db := openDb()
	defer db.Close()

	sql := `select Movie, Name, Offset, Content, Type from subtitle where Name=?`
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
	db := openDb()
	defer db.Close()

	var count int
	err := db.QueryRow("select count(*) from subtitle where Name=?", sub.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		sql := "insert into subtitle(Movie, Name, Offset, Content, Type) values (?,?,?,?,?)"
		_, err := db.Exec(sql, sub.Movie, sub.Name, sub.Offset, sub.Content, sub.Type)
		if err != nil {
			log.Print(err)
		}
	}
}

func UpdateSubtitleOffset(name string, offset time.Duration) {
	db := openDb()
	defer db.Close()

	_, err := db.Exec("update subtitle set Offset=? where Name=?", int64(offset), name)
	if err != nil {
		log.Print(err)
	}
}

func scanPlaying(scanner rowScanner) (*Playing, error) {
	var p Playing
	var lastPos int64
	err := scanner.Scan(&p.Movie, &lastPos, &p.SoundStream, &p.Sub1, &p.Sub2)
	if err == nil {
		p.LastPos = time.Duration(lastPos)
		return &p, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetPlaying(movie string) *Playing {
	db := openDb()
	defer db.Close()

	sql := `select Movie, LastPos, SoundStream, Sub1, Sub2 from Playing where Movie=?`
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
		SavePlaying(&Playing{movie, 0, -1, "", ""})
	}

	return GetPlaying(movie)
}

func SavePlaying(p *Playing) {
	log.Printf("Save playing: %v", *p)

	db := openDb()
	defer db.Close()

	var count int
	err := db.QueryRow("select count(*) from playing where Movie=?", p.Movie).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		sql := "insert into playing(Movie, LastPos, SoundStream, Sub1, Sub2) values (?,?,?,?,?)"
		_, err := db.Exec(sql, p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := db.Exec("update playing set Movie=?,LastPos=?,SoundStream=?,Sub1=?,Sub2=? where Movie=?",
			p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2, p.Movie)
		if err != nil {
			log.Print(err)
		}
	}
}
