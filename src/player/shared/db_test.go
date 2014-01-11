package shared

import (
	// "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

func setup() {
	DbFile = "./test.db"

	db := openDb()
	defer db.Close()

	_, err := db.Exec("delete from subtitle")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("delete from playing")
	if err != nil {
		log.Fatal(err)
	}

	//	drop table playing; create table playing (Movie nvarchar(2048) not null primary key default(''), LastPos int not null default(0), SoundStream int not null default(-1), Sub1 nvarchar(2048) not null default(''), Sub2 nvarchar(2048) not null default(''));
	// drop table subtitle; create table subtitle (Name nvarchar(2048) not null default(''), Movie nvarchar(2048) not null default(''), Offset int not null default(0), Content text not null default(''));

	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestInsertSubtitle(t *testing.T) {
	setup()

	InsertSubtitle(&Sub{"moviename", "subname", 1, "subcontent"})

	subs := GetSubtitles("moviename")
	if len(subs) != 1 {
		t.Errorf("subs length: expect %d, got %d", 1, len(subs))
		return
	}

	if subs[0].Name != "subname" {
		t.Errorf("sub name expect 'subname', got '%s'", subs[0].Name)
	}

	if int64(subs[0].Offset) != 1 {
		t.Errorf("sub offset expect 1, got %d", int64(subs[0].Offset))
	}

	if subs[0].Content != "subcontent" {
		t.Errorf("sub content expect 'subcontent', got '%s'", subs[0].Content)
	}
}

func TestUpdateSubtitle(t *testing.T) {
	setup()

	InsertSubtitle(&Sub{"moviename", "subname", 1, "subcontent"})

	UpdateSubtitleOffset("subname", 100)

	subs := GetSubtitles("moviename")
	if len(subs) != 1 {
		t.Errorf("subs length: expect %d, got %d", 1, len(subs))
		return
	}

	if subs[0].Name != "subname" {
		t.Errorf("sub name expect 'subname', got '%s'", subs[0].Name)
	}

	if int64(subs[0].Offset) != 100 {
		t.Errorf("sub offset expect 100, got %d", int64(subs[0].Offset))
	}

	if subs[0].Content != "subcontent" {
		t.Errorf("sub content expect 'subcontent', got '%s'", subs[0].Content)
	}
}

func TestPlaying(t *testing.T) {
	setup()

	p := &Playing{"moviename", 123, 8, "sub1name", "sub2name"}
	SavePlaying(p)

	p1 := GetPlaying("moviename")

	if p.LastPos != p1.LastPos {
		t.Errorf("p.LastPost: expect %s, got %s", p.LastPos, p1.LastPos)
	}

	if p.SoundStream != p1.SoundStream {
		t.Errorf("p.SoundStream: expect %s, got %s", p.SoundStream, p1.SoundStream)
	}

	if p.Sub1 != p1.Sub1 {
		t.Errorf("p.Sub1: expect %s, got %s", p.Sub1, p1.Sub1)
	}
	if p.Sub2 != p1.Sub2 {
		t.Errorf("p.Sub2: expect %s, got %s", p.Sub2, p1.Sub2)
	}
}
