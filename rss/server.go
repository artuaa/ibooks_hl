package rss

import (
	"highlights/ibooks"
	"io"
	"log"
	"net/http"
	"strconv"
)

var defaultCount = 3

type Web struct {
	rss RSS
}

func NewWeb() *Web {
	storage := ibooks.NewStorage()
	return &Web{rss: New(storage)}
}

func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rss/highlights" {
		log.Printf("request %s", r.URL.Path)
		w.Header().Set("content-type", "application/rss+xml")
		var count = defaultCount
		queryParams := r.URL.Query()
		if value, ok := queryParams["count"]; ok {
			if v, err := strconv.Atoi(value[0]); err == nil {
				count = v
			}
		}
		rss, err := web.rss.GenerateFeed(count)
		if err != nil {
			log.Printf("generate feed error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, rss)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func RunServer() error {
	w := NewWeb()
	port := "9999"
	log.Printf("server started on port: %s", port)
	err := http.ListenAndServe(":"+port, w)
	if err != nil {
		return err
	}
	return nil
}
