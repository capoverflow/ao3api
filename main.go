package ao3api

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/capoverflow/ao3api/internal/scrapper"
	"github.com/capoverflow/ao3api/internal/utils"
	"github.com/capoverflow/ao3api/models"
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

// Tags
func Search(SearchString models.Search) {
	scrapper.Search(SearchString)

}

func Users(Params models.UserParams) (AuthorInfo models.User) {

	log.Println(Params)

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
	Booksmarks := UserBookmarks(Params)
	Pseuds := strings.Split(AuthorInfo.Profile.Pseuds, ",")
	log.Println(len(Pseuds))
	var WorksID []string
	for _, pseud := range Pseuds {

		Params.Pseuds = strings.TrimSpace(pseud)
		Works := scrapper.GetUsersWorks(Params)
		AuthorInfo.Username = Params.Username
		for _, work := range Works {
			// log.Println(work)
			splitWorks := strings.Split(work, "/")
			// log.Println(splitWorks[len(splitWorks)-1], len(splitWorks))
			if len(splitWorks) != 0 {
				WorksID = append(WorksID, splitWorks[len(splitWorks)-1])
			}

		}
	}

	Works := scrapper.GetUsersWorks(Params)
	AuthorInfo.Username = Params.Username
	for _, work := range Works {
		// log.Println(work)
		splitWorks := strings.Split(work, "/")
		// log.Println(splitWorks[len(splitWorks)-1], len(splitWorks))
		if len(splitWorks) != 0 {
			WorksID = append(WorksID, splitWorks[len(splitWorks)-1])
		}

	}

	if len(Pseuds) > 1 {
		DedupWorks := utils.RemoveDuplicates(WorksID)

		// fmt.Printf("Len before dedup %d\nLen after dedup %d\n", len(WorksID), len(DedupWorks))

		for _, work := range DedupWorks {
			AuthorInfo.Works = append(AuthorInfo.Works, models.Work{
				WorkID: work,
			})
		}
	} else {
		for _, work := range WorksID {
			AuthorInfo.Works = append(AuthorInfo.Works, models.Work{
				WorkID: work,
			})
		}
	}

	for _, work := range Booksmarks {
		s := strings.Replace(work, "/works/", "", -1)
		AuthorInfo.Bookmarks = append(AuthorInfo.Bookmarks, models.Work{
			WorkID: s,
		})
	}

	return AuthorInfo
}

func UserBookmarks(Params models.UserParams) (Bookmarks []string) {

	if len(Params.Addr) == 0 {
		if Params.Debug {
			log.Println("Debug: Setting default url")
		}
		Params.Addr = "archiveofourown.org"
	}
	if Params.Debug {
		log.Printf("Debug:\n %v ", Params)
	}

	Bookmarks = scrapper.GetUserBookmarks(Params)

	return Bookmarks
}
