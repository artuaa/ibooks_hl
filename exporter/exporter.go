package exporter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"highlights/ibooks"

	_ "github.com/mattn/go-sqlite3"
)

func saveToFile(path string, hls []ibooks.Highlight) error {
	_, err := os.Stat(path)
	if err == nil {
		log.Printf("Skip book notes '%s'. File already exists", path)
		return nil
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	for i, v := range hls {
		if i == 0 || v.Chapter.String != hls[i-1].Chapter.String && v.Chapter.String != "" {
			f.WriteString(fmt.Sprintf("### %s\n", v.Chapter.String))
		}
		if v.Text.Valid {
			f.WriteString(fmt.Sprintf("> %s\n%s\n", strings.Replace(v.Text.String, "\n", "\n> ", -1), v.Note.String))
		}
	}
	//if err := f.Close(); err != nil {
	//	return err
	//}
	return nil
}

// ExportNotes Export ibooks notes and highlights to obsidian vault
func ExportNotes() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify notes path")
	}
	notesPath := os.Args[1]

	cmd := exec.Command("mkdir", "-p", notesPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Can't create target directory '%s' %s", notesPath, err)
	}
	storage := ibooks.Storage{}
	highlights, err := storage.LoadHighlights()
	if err != nil {
		log.Fatal(err)
	}
	for title, hls := range ibooks.GroupByTitle(highlights) {
		if len(hls) < 12 {
			continue
		}
		path := notesPath + "/" + title + ".md"
		err := saveToFile(path, hls)
		if err != nil {
			log.Fatal(err)
		}
	}
}
