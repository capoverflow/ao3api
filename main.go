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
func Fanfic(Params models.FanficParams) (fanfic models.Work, status int, err error) {
	var ChaptersIDs []string
	// var fanfic ao3structs.Work

	ChaptersIDs, status, err = scrapper.GetFirstChapterID(Params.WorkID, Params.ChapterID, Params.ProxyURLs, Params.Debug)
	// if len(proxyURLs) != 0 {
	// 	log.Println(proxyURLs)
	// }
	// log.Println("ChaptersIDs: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs))
	if status != 404 {
		if len(ChaptersIDs) != 0 {
			fanfic = scrapper.GetInfo(Params.WorkID, ChaptersIDs, Params.ProxyURLs)
			fanfic.WorkID = Params.WorkID
			fanfic.ChapterID = Params.ChapterID
			fanfic.URL = fmt.Sprintf("https://archiveofourown.org/works/%s", Params.WorkID)
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

func UserBookmarks(Author string) (Bookmarks []string) {
	log.Println(Author)

	// Bookmarks = utils.RemoveDuplicateStr(scrapper.GetUserBookmarks(Author))
	Bookmarks = scrapper.GetUserBookmarks(Author)

	log.Println(len(Bookmarks))

	return Bookmarks
}
