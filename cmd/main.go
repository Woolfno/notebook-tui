package main

import (
	"log"
	"notebook/internal/tui"
)

func main() {
	notesDir := "notes"
	t := tui.New(notesDir)
	if err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
