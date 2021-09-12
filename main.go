package ao3

import (
	"ao3api/internal/models"
	"ao3api/internal/scrapper"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("ao3api: ")
	log.SetOutput(os.Stderr)
}

// Parsing parse the fanfiction from ao3
func Fanfic(WorkID, ChapterID string, debug bool) (fanfic models.Work, status int, err error) {
	var ChaptersIDs []string
	// var fanfic ao3structs.Work

	ChaptersIDs, status, err = scrapper.GetFirstChapterID(WorkID, ChapterID, debug)
	// log.Println("ChaptersIDs: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs))
	if status != 404 {
		if len(ChaptersIDs) != 0 {
			fanfic = scrapper.GetInfo(WorkID, ChaptersIDs)
			fanfic.WorkID = WorkID
			fanfic.ChapterID = ChapterID
		}
	} else {
		log.Println("status 404")
	}
	// log.Println(WorkID, ChapterID, status)
	return fanfic, status, err
}

//Tags
func Search(search string) {
	// log.Println(scrapper.Search(search))
	scrapper.Search(search)
	// scrapper.Tags()
}
