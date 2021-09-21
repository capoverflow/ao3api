package ao3

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/capoverflow/ao3api/internal/scrapper"
	"gitlab.com/capoverflow/ao3api/models"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("ao3api: ")
	log.SetOutput(os.Stderr)
}

// Parsing parse the fanfiction from ao3
func Fanfic(WorkID, ChapterID string, debug bool, proxyURLs []string) (fanfic models.Work, status int, err error) {
	var ChaptersIDs []string
	// var fanfic ao3structs.Work

	ChaptersIDs, status, err = scrapper.GetFirstChapterID(WorkID, ChapterID, proxyURLs, debug)
	// if len(proxyURLs) != 0 {
	// 	log.Println(proxyURLs)
	// }
	// log.Println("ChaptersIDs: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs))
	if status != 404 {
		if len(ChaptersIDs) != 0 {
			fanfic = scrapper.GetInfo(WorkID, ChaptersIDs, proxyURLs)
			fanfic.WorkID = WorkID
			fanfic.ChapterID = ChapterID
			fanfic.URL = fmt.Sprintf("https://archiveofourown.org/works/%s", WorkID)
		}
	} else {
		log.Println("status 404")
	}
	// log.Println(WorkID, ChapterID, status)
	return fanfic, status, err
}

//Tags
func Search(SearchString models.Search) {
	scrapper.Search(SearchString)

}

func Users(Author string) (AuthorInfo models.User) {
	AuthorInfo = scrapper.GetUsersInfo(Author)
	return AuthorInfo
}
