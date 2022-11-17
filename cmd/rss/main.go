package main

import "highlights/rss"

func main() {
	err := rss.RunServer()
	if err != nil {
		panic(err)
	}
}
