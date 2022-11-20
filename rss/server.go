package rss

import (
	"fmt"
	"highlights/ibooks"
	"io"
	"net/http"
)

func RunServer() error {
	http.HandleFunc("/", root)
	port := "9999"
	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}

//routes
func root(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/rss+xml")
	r := New(&ibooks.IBooksStorage{})
	rss, err := r.GenerateFeed(3)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, rss)
	}
}
