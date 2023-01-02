package rss

import (
	"fmt"
	"highlights/ibooks"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

type RSS struct {
	storage *ibooks.Storage
}

func New(s *ibooks.Storage) RSS {
	return RSS{storage: s}
}

func makeObsidianUrl(title string) string {
	v := url.Values{}
	v.Add("file", "book_highlights/"+title+".md")
	v.Add("vault", "vault")
	return "obsidian://open?" + v.Encode()
}

func wrapSentence(s string) string {
	if strings.HasSuffix(s, ".") {
		return s
	} else {
		return s + "..."
	}
}

func (r *RSS) randomNotes(count int, at time.Time) ([]*feeds.Item, error) {
	hls, err := r.storage.LoadHighlights()
	if err != nil {
		return nil, fmt.Errorf("can't load highlights %w", err)
	}
	var result []*feeds.Item
	// HACK:
	if len(hls) < count {
		count = len(hls)
	}
	rand.Seed(at.Unix())
	for i := 0; i < count; i++ {
		idx := rand.Intn(len(hls))
		hl := hls[idx]

		item := &feeds.Item{
			Title:       hl.Title.String,
			Description: wrapSentence(hl.Text.String),
			Link:        &feeds.Link{Href: makeObsidianUrl(hl.Title.String)},
			Author:      &feeds.Author{Name: hl.Author.String},
			Created:     at.Add(time.Second * time.Duration(i)),
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *RSS) GenerateFeed(count int) (string, error) {
	now := time.Now()
	createdAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	feed := &feeds.Feed{
		Title:       "Books highlights",
		Link:        &feeds.Link{Href: "http://google.com"},
		Description: "quotes and words",
		Author:      &feeds.Author{Name: "Artur Sharipov", Email: "theartua@gmail.com"},
		Created:     createdAt,
	}

	items, err := r.randomNotes(count, createdAt)

	if err != nil {
		return "", fmt.Errorf("can't load random notes. count %d, date %v %w", count, createdAt, err)
	}
	feed.Items = items
	rss, err := feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("can't generate rss feed %w", err)
	}

	return rss, nil
}
