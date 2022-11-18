package rss

import (
	"highlights/ibooks"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/feeds"
)

type RSS struct {
	storage ibooks.Storage
}

func New(s ibooks.Storage) RSS {
	return RSS{storage: s}
}

func makeObsidianUrl(title string) string {
	v := url.Values{}
	v.Add("file", "book_highligths/"+title+".md")
	v.Add("vault", "vault")
	return "obsidian://open?" + v.Encode()
}

func (r *RSS) randomNotes(count int) ([]*feeds.Item, error) {
	hls, err := r.storage.LoadHighlights()
	if err != nil {
		return nil, err
	}
	result := []*feeds.Item{}
	// HACK:
	if len(hls) < count {
		count = len(hls)
	}
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	for i := 0; i < count; i++ {
		idx := rand.Intn(len(hls))
		hl := hls[idx]

		item := &feeds.Item{
			Title:       hl.Title.String,
			Description: hl.Text.String,
			Link:        &feeds.Link{Href: makeObsidianUrl(hl.Title.String)},
			Author:      &feeds.Author{Name: hl.Author.String},
			Created:     now.Add(-time.Second * time.Duration(i)),
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *RSS) GenerateFeed(count int) (string, error) {
	now := time.Now()

	feed := &feeds.Feed{
		Title: "Artua books highlights",
		// TODO: add link
		Link:        &feeds.Link{Href: "http://google.com"},
		Description: "quotes and words",
		Author:      &feeds.Author{Name: "Artur Sharipov", Email: "theartua@gmail.com"},
		Created:     now,
	}

	items, err := r.randomNotes(count)

	if err != nil {
		return "", err
	}
	feed.Items = items
	rss, err := feed.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	return rss, nil
}
