package ao3

import (
	"log"
	"os"

	"gitlab.com/capoverflow/ao3api/internal/ao3structs"
	"gitlab.com/capoverflow/ao3api/internal/scrapper"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("ao3api: ")
	log.SetOutput(os.Stderr)
}

// Parsing parse the fanfiction from ao3
func Fanfic(WorkID, ChapterID string, debug bool) (ao3structs.Work, int) {

	var ChaptersIDs []string
	var status int
	var fanfic ao3structs.Work

	ChaptersIDs, status = scrapper.GetFirstChapterID(WorkID, ChapterID, debug)
	// log.Println("ChaptersID: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs), err)
	if status != 404 {
		fanfic = scrapper.GetInfo(WorkID, ChaptersIDs)
		fanfic.WorkID = WorkID
		fanfic.ChapterID = ChapterID

	} //else {
	// 	log.Println("status 404")
	// }
	// log.Println(WorkID, ChapterID, status)
	return fanfic, status
}

//Tags
func Search(search string) {
	// log.Println(scrapper.Search(search))
	scrapper.Search(search)
	// scrapper.Tags()
}
