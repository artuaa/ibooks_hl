package rss

import (
	"database/sql"
	"highlights/ibooks"
	"testing"
)

// type Feed struct {
// 	Version     string `json:"version"`
// 	Title       string `json:"title"`
// 	HomePageURL string `json:"home_page_url"`
// 	Description string `json:"description"`
// 	Author      struct {
// 		Name string `json:"name"`
// 	} `json:"author"`
// 	Items []struct {
// 		ID            string    `json:"id"`
// 		URL           string    `json:"url"`
// 		Title         string    `json:"title"`
// 		Summary       string    `json:"summary"`
// 		DatePublished time.Time `json:"date_published"`
// 		Author        struct {
// 			Name string `json:"name"`
// 		} `json:"author"`
// 	} `json:"items"`
// }

// var expected []Feed = []Feed{{
// 	Version:     "",
// 	Title:       "",
// 	HomePageURL: "",
// 	Description: "",
// 	Author:      struct{Name string "json:\"name\""}{},
// 	Items:       []struct{ID string "json:\"id\""; URL string "json:\"url\""; Title string "json:\"title\""; Summary string "json:\"summary\""; DatePublished time.Time "json:\"date_published\""; Author struct{Name string "json:\"name\""} "json:\"author\""}{},
// }}

// {
// 	"version": "https://jsonfeed.org/version/1",
// 	"title": "Artua books highlights",
// 	"home_page_url": "http://google.com",
// 	"description": "quotes and words",
// 	"author": {
// 	  "name": "Artur Sharipov"
// 	},
// 	"items": [
// 	  {
// 		"id": "",
// 		"url": "obsidian://open?file=book_highligths%2FTitle.md\u0026vault=vault",
// 		"title": "Title",
// 		"summary": "Text",
// 		"date_published": "2022-11-18T10:54:38.337773+03:00",
// 		"author": {
// 		  "name": "Author"
// 		}
// 	  }
// 	]
//   }, want

type MockStorage struct{}

func (s *MockStorage) LoadHighlights() ([]ibooks.Highlight, error) {
	hls := []ibooks.Highlight{
		{
			Author:  sql.NullString{String: "Author"},
			Title:   sql.NullString{String: "Title"},
			Chapter: sql.NullString{String: "Chapter"},
			Text:    sql.NullString{String: "Text"},
			Note:    sql.NullString{String: "Note"},
		},
	}
	return hls, nil
}

func TestGenerateFeed(t *testing.T) {
	storage := &MockStorage{}
	rss := New(storage)
	tests := []struct {
		name    string
		count   int
		want    string
		wantErr bool
	}{
		{"base test", 5, "", false},
		{"base test", 3000, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rss.GenerateFeed(4)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
