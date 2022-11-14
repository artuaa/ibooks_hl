package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Highlight struct {
	Note    sql.NullString `db:"note"`
	Author  sql.NullString `db:"author"`
	Title   sql.NullString `db:"title"`
	Text    sql.NullString `db:"selected_text"`
	Chapter sql.NullString `db:"chapter"`
}

func loadHighlights() ([]Highlight, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	annotations_path := homedir + "/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation/AEAnnotation_v10312011_1727_local.sqlite"
	library_path := homedir + "/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary/BKLibrary-1-091020131601.sqlite"

	query := `select
	--ZANNOTATIONASSETID as asset_id,
	ZTITLE as title,
	ZAUTHOR as author,
	ZANNOTATIONSELECTEDTEXT as selected_text,
	ZANNOTATIONNOTE as note,
	--ZANNOTATIONREPRESENTATIVETEXT as represent_text,
	ZFUTUREPROOFING5 as chapter
	--ZANNOTATIONSTYLE as style,
	--ZANNOTATIONMODIFICATIONDATE as modified_date,
	--ZANNOTATIONLOCATION as location
	from ZAEANNOTATION
	left join books.ZBKLIBRARYASSET
	on ZAEANNOTATION.ZANNOTATIONASSETID = books.ZBKLIBRARYASSET.ZASSETID
	order by ZANNOTATIONASSETID, ZPLLOCATIONRANGESTART;`

	db, err := sqlx.Connect("sqlite3", annotations_path)
	if err != nil {
		return nil, fmt.Errorf("open db error: %s", err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("attach database '%s' as books;", library_path))
	if err != nil {
		return nil, err
	}
	rows, err := db.Queryx(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := []Highlight{}
	for rows.Next() {
		hl := Highlight{}
		err = rows.StructScan(&hl)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, hl)
	}
	return result, nil
}
func groupByTitle(hls []Highlight) map[string][]Highlight {
	result := make(map[string][]Highlight)
	for _, v := range hls {
		if _, ok := result[v.Title.String]; !ok {
			result[v.Title.String] = []Highlight{}
		}
		result[v.Title.String] = append(result[v.Title.String], v)
	}
	return result
}

func saveToFile(path string, hls []Highlight) error {
	f, err := os.OpenFile(path+".md", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	for _, hls := range groupByTitle(hls) {
		f.WriteString(fmt.Sprintf("# %s\n%s\n", hls[0].Title.String, hls[0].Author.String))
		for _, v := range hls {
			if v.Text.Valid {
				f.WriteString(fmt.Sprintf("### %s\n> %s\n%s\n", v.Chapter.String, strings.Replace(v.Text.String, "\n", "\n> ", -1), v.Note.String))
			}
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify notes path")
		return
	}
	notes_path := os.Args[1]
	highlights, err := loadHighlights()
	if err != nil {
		log.Fatal(err)
	}
	for title, hls := range groupByTitle(highlights) {
		if len(hls) < 12 {
			continue
		}
		err := saveToFile(notes_path+"/"+title, hls)
		if err != nil {
			log.Fatal(err)
		}
	}
}
