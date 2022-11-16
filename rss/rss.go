package rss

import (
	"highligths/ibooks"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/feeds"
)

func randomNotes(count int) ([]*feeds.Item, error) {
	hls, err := ibooks.LoadHighlights()
	if err != nil {
		return nil, err
	}
	rand.Seed(499)
	result := []*feeds.Item{}
	// HACK:
	if len(hls) < count {
		count = len(hls)
	}
	for i := 0; i < count; i++ {
		idx := rand.Intn(len(hls))
		hl := hls[idx]
		now := time.Now()
		item := &feeds.Item{
			Title:       hl.Title.String,
			Description: hl.Text.String,
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		}
		result = append(result, item)
	}
	return result, nil
}

func GenerateFeed(count int) (string, error) {
	now := time.Now()
	feed := &feeds.Feed{
		Title: "Artua books highlights",
		// TODO: add link
		Link:        &feeds.Link{Href: "http://google.com"},
		Description: "quotes and words",
		Author:      &feeds.Author{Name: "Artur Sharipov", Email: "theartua@gmail.com"},
		Created:     now,
	}

	items, err := randomNotes(count)

	if err != nil {
		return "", err
	}
	feed.Items = items
	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	return rss, nil
}
