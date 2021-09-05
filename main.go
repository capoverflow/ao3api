package ao3

import (
	"log"

	"gitlab.com/capoverflow/ao3api/internal/ao3structs"
	"gitlab.com/capoverflow/ao3api/internal/scrapper"
)

//package main

// Work is the struc with the fanfic info

// Parsing parse the fanfiction from ao3
func Fanfic(WorkID, ChapterID string, debug bool) (ao3structs.Work, int) {
	// log.Panic("parsing test")

	var ChaptersIDs []string
	var status int
	var fanfic ao3structs.Work

	ChaptersIDs, status = scrapper.GetFirstChapterID(WorkID, ChapterID, debug)
	// log.Println("ChaptersID: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs), err)
	if status != 404 {
		fanfic = scrapper.GetInfo(WorkID, ChaptersIDs)
		fanfic.WorkID = WorkID
		fanfic.ChapterID = ChapterID

	} else {
		log.Panic("status 404")
	}
	// log.Println(WorkID, ChapterID, status)
	return fanfic, status
}

//Tags
func Tags(tags string) {
	log.Println(scrapper.Tags(tags))
	// scrapper.Tags()
}
