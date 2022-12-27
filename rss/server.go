package rss

import (
	"fmt"
	"highlights/ibooks"
	"io"
	"net/http"
)

type Web struct {
	rss RSS
}

func NewWeb() *Web {
	storage := ibooks.NewStorage()
	return &Web{rss: New(storage)}
}

func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rss/highlights" {
		w.Header().Set("content-type", "application/rss+xml")
		rss, err := web.rss.GenerateFeed(3)
		if err != nil {
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
	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":"+port, w)
	if err != nil {
		return err
	}
	return nil
}
