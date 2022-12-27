package ibooks

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	annotationsPath = "/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation/AEAnnotation_v10312011_1727_local.sqlite"
	libraryPath     = "/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary/BKLibrary-1-091020131601.sqlite"
)

type Highlight struct {
	Note    sql.NullString `db:"note"`
	Author  sql.NullString `db:"author"`
	Title   sql.NullString `db:"title"`
	Text    sql.NullString `db:"selected_text"`
	Chapter sql.NullString `db:"chapter"`
}

type Storage struct {
	annotationsPath string
	libraryPath     string
}

func NewStorage() *Storage {
	return &Storage{
		annotationsPath: annotationsPath,
		libraryPath:     libraryPath,
	}
}

func (s *Storage) LoadHighlights() ([]Highlight, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
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
	where selected_text is not null
	order by ZANNOTATIONASSETID, ZPLLOCATIONRANGESTART;`

	db, err := sqlx.Connect("sqlite3", homedir+s.annotationsPath)
	defer db.Close()
	if err != nil {
		return nil, fmt.Errorf("open db error: %s", err)
	}
	_, err = db.Exec(fmt.Sprintf("attach database '%s' as books;", homedir+s.libraryPath))
	if err != nil {
		return nil, err
	}
	rows, err := db.Queryx(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var result []Highlight
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

func GroupByTitle(hls []Highlight) map[string][]Highlight {
	result := make(map[string][]Highlight)
	for _, v := range hls {
		if _, ok := result[v.Title.String]; !ok {
			result[v.Title.String] = []Highlight{}
		}
		result[v.Title.String] = append(result[v.Title.String], v)
	}
	return result
}
