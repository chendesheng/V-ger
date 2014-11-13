package shared

import (
	"log"
	"time"
	"vger/dbHelper"

	_ "github.com/mattn/go-sqlite3"
)

func scanSub(scanner dbHelper.RowScanner) (*Sub, error) {
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
	db := dbHelper.Open()
	defer dbHelper.Close(db)

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

func GetSubtitlesMap(movie string) map[string]*Sub {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	sql := `select Movie, Name, Offset, Content, Type, Lang1, Lang2 from subtitle where Movie=?`
	rows, err := db.Query(sql, movie)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	subs := make(map[string]*Sub)
	for rows.Next() {
		sub, err := scanSub(rows)
		if err == nil {
			subs[sub.Name] = sub
		}
	}

	return subs
}

func GetSubtitle(name string) *Sub {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

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
	db := dbHelper.Open()
	defer dbHelper.Close(db)

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

func DeleteSubtitle(movie string) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("delete from subtitle where Movie=?", movie)
	return err
}

func DeletePlaying(movie string) error {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("delete from playing where Movie=?", movie)
	return err
}

var chUpdateSubtitleOffset chan map[string]interface{}

func UpdateSubtitleOffsetAsync(name string, offset time.Duration) {
	if chUpdateSubtitleOffset == nil {
		chUpdateSubtitleOffset = make(chan map[string]interface{}, 20)
		go func() {
			timer := time.NewTicker(2 * time.Second)
			var arg map[string]interface{}
			var lastArg map[string]interface{}
			for {
				select {
				case arg = <-chUpdateSubtitleOffset:
					break
				case <-timer.C:
					if len(arg) > 0 &&
						(len(lastArg) == 0 || arg["name"].(string) != lastArg["name"].(string) ||
							arg["offset"].(time.Duration) != lastArg["offset"].(time.Duration)) {
						UpdateSubtitleOffset(arg["name"].(string), arg["offset"].(time.Duration))
						lastArg = arg
					}
					break
				}
			}
		}()
	}
	chUpdateSubtitleOffset <- map[string]interface{}{
		"name":   name,
		"offset": offset,
	}
}

func UpdateSubtitleOffset(name string, offset time.Duration) {
	log.Printf("update subtitle offset: %s, %d", name, offset)
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("update subtitle set Offset=? where Name=?", int64(offset), name)
	if err != nil {
		log.Print(err)
	}
}

func UpdateSubtitleLanguage(name string, lang1, lang2 string) {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	_, err := db.Exec("update subtitle set Lang1=?,Lang2=? where Name=?", lang1, lang2, name)
	if err != nil {
		log.Print(err)
	}
}

func scanPlaying(scanner dbHelper.RowScanner) (*Playing, error) {
	var p Playing
	var lastPos, duration int64
	err := scanner.Scan(&p.Movie, &lastPos, &p.SoundStream, &p.Sub1, &p.Sub2, &duration, &p.Volume)
	if err == nil {
		p.LastPos = lastPos
		p.Duration = time.Duration(duration)
		p.FirstOpen = lastPos == 0
		return &p, nil
	} else {
		log.Print(err)
		return nil, err
	}
}

func GetPlaying(movie string) *Playing {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	sql := `select Movie, LastPos, SoundStream, Sub1, Sub2, Duration, Volume from Playing where Movie=?`
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

func PlayingExists(movie string) bool {
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	var count int
	err := db.QueryRow("select count(*) from playing where Movie=?", movie).Scan(&count)
	if err != nil {
		log.Print(err)
	}
	return count > 0
}

func CreateOrGetPlaying(movie string) *Playing {
	if !PlayingExists(movie) {
		SavePlaying(&Playing{movie, 0, -1, "", "", 0, 8, 0, true})
	}

	return GetPlaying(movie)
}

var chPlaying chan *Playing

func SavePlayingAsync(p *Playing) {
	if chPlaying == nil {
		chPlaying = make(chan *Playing, 20)
		go func() {
			timer := time.NewTicker(2 * time.Second)
			var p *Playing
			// var lastP *Playing
			for {
				select {
				case p = <-chPlaying:
					break
				case <-timer.C:
					// if p != nil && p != lastP {
					SavePlaying(p)
					// lastP = p
					// }
					break
				}
			}
		}()
	}

	chPlaying <- p
}
func SavePlaying(p *Playing) {
	if p == nil {
		return
	}

	// log.Printf("Save playing: %v", p)
	db := dbHelper.Open()
	defer dbHelper.Close(db)

	var count int
	err := db.QueryRow("select count(*) from playing where Movie=?", p.Movie).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		sql := "insert into playing(Movie, LastPos, SoundStream, Sub1, Sub2, Duration, Volume) values (?,?,?,?,?,?,?)"
		_, err := db.Exec(sql, p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2, p.Duration, p.Volume)
		if err != nil {
			log.Print(err)
		}
	} else {
		_, err := db.Exec("update playing set Movie=?,LastPos=?,SoundStream=?,Sub1=?,Sub2=?,Duration=?,Volume=? where Movie=?",
			p.Movie, p.LastPos, p.SoundStream, p.Sub1, p.Sub2, p.Duration, p.Volume, p.Movie)
		if err != nil {
			log.Print(err)
		}
	}
}
