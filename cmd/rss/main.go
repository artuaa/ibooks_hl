package main

import "highligths/rss"

func main() {
	err := rss.RunServer()
	if err != nil {
		panic(err)
	}
}
