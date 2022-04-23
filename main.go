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
	if len(Params.Addr) == 0 {
		if Params.Debug {
			log.Println("Debug: Setting default url")
		}
		Params.Addr = "archiveofourown.org"
	}
	if Params.Debug {
		log.Printf("Debug:\n %v ", Params)
	}
	var ChaptersIDs []string
	// var fanfic ao3structs.Work

	ChaptersIDs, status, err = scrapper.GetFirstChapterID(Params)

	if status != 404 {
		if len(ChaptersIDs) != 0 {
			fanfic = scrapper.GetInfo(Params, ChaptersIDs)
			fanfic.WorkID = Params.WorkID
			fanfic.ChapterID = Params.ChapterID
			fanfic.URL = fmt.Sprintf("http://%s/works/%s", Params.Addr, Params.WorkID)
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

func Users(Params models.UserParams) (AuthorInfo models.User) {
	if len(Params.Addr) == 0 {
		if Params.Debug {
			log.Println("Debug: Setting default url")
		}
		Params.Addr = "archiveofourown.org"
	}
	if Params.Debug {
		log.Printf("Debug:\n %v ", Params)
	}

	AuthorInfo = scrapper.GetUsersInfo(Params)
	return AuthorInfo
}

func UserBookmarks(Author string) (Bookmarks []string) {
	log.Println(Author)

	// Bookmarks = utils.RemoveDuplicateStr(scrapper.GetUserBookmarks(Author))
	Bookmarks = scrapper.GetUserBookmarks(Author)

	log.Println(len(Bookmarks))

	return Bookmarks
}
