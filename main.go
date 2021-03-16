package ao3

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

//package main

// Work is the struc with the fanfic info
type Work struct {
	URL             string   `json:"URL,omitempty"`
	WorkID          string   `json:"WorkID,omitempty"`
	ChapterID       string   `json:"ChapterID,omitempty"`
	ChapterTitle    string   `json:"ChapterTitle,omitempty"`
	Title           string   `json:"Title,omitempty"`
	Author          string   `json:"Author,omitempty"`
	Published       string   `json:"Published,omitempty"`
	Updated         string   `json:"Updated,omitempty"`
	Words           string   `json:"Words,omitempty"`
	Chapters        string   `json:"Chapters,omitempty"`
	Comments        string   `json:"Comments,omitempty"`
	Kudos           string   `json:"Kudos,omitempty"`
	Bookmarks       string   `json:"Bookmarks,omitempty"`
	Hits            string   `json:"Hits,omitempty"`
	Fandom          string   `json:"Fandom,omitempty"`
	Summary         []string `json:"Summary,omitempty"`
	ChaptersTitles  []string `json:"ChaptersTitles,omitempty"`
	ChaptersIDs     []string `json:"ChaptersIDs,omitempty"`
	Relationship    []string `json:"Relationship,omitempty"`
	AlternativeTags []string `json:"AlternativeTags,omitempty"`
}

type id struct {
	WorkID    string
	ChapterID string
}

type ids struct {
	works []id
}

type fanfic struct {
	WorkID    string
	ChapterID string
	Debug     bool
}

// Parsing parse the fanfiction from ao3
func Parsing(WorkID, ChapterID string, debug bool) Work {
	var ChaptersIDs []string
	var err bool = false
	ChaptersIDs, err = getFirstChapterID(WorkID, ChapterID, debug)
	log.Println("ChaptersID: , ChaptersIDs length:", ChaptersIDs, len(ChaptersIDs))
	fanfic := getInfo(WorkID, ChaptersIDs)
	if err != false {
		panic(err)
	}
	fanfic.WorkID = WorkID
	fanfic.ChapterID = ChapterID
	//log.Println("getInfo: ", fanfic)
	return fanfic

}
func getFirstChapterID(WorkID, ChapterID string, debug bool) ([]string, bool) {
	var err = false

	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", WorkID)
	log.Printf("WorkID: %s, url %s", WorkID, url)
	//url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate", wID)
	// var title string
	// var author string
	var ChaptersIDs []string
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*archiveofourown.org/*",
		// Set a delay between requests to these domains
		Delay: 5 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
		// Add User Agent
		Parallelism: 2,
	})
	c.OnHTML("#main > ol", func(e *colly.HTMLElement) {
		hrefChaptersIDs := e.ChildAttrs("a", "href")
		ChaptersIDs = findChaptersIDs(hrefChaptersIDs)
	})

	c.OnRequest(func(r *colly.Request) {
		if debug == true {
			log.Println("visiting", r.URL.String())
		}

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Println(r.Request)
			log.Println(string(r.Body))
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		// StatusCode := r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		// StatusCode := r.StatusCode
	})

	c.Visit(url)
	return ChaptersIDs, err
}

func getInfo(WorkID string, ChaptersIDs []string) Work {
	var Fanfic Work
	Fanfic.ChaptersIDs = ChaptersIDs
	Fanfic.URL = fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s?view_adult=true", WorkID, Fanfic.ChaptersIDs[0])
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		//colly.AllowURLRevisit(),
	)
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "*archiveofourown.org/*",
		// Set a delay between requests to these domains
		Delay: 5 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
		// Add User Agent
		Parallelism: 2,
	})

	c.OnHTML("dl.stats", func(e *colly.HTMLElement) {
		Fanfic.Published = e.ChildText("dd.published")
		Fanfic.Updated = e.ChildText("dd.status")
		Fanfic.Words = e.ChildText("dd.words")
		Fanfic.Chapters = e.ChildText("dd.chapters")
		Fanfic.Comments = e.ChildText("dd.comments")
		Fanfic.Kudos = e.ChildText("dd.kudos")
		Fanfic.Bookmarks = e.ChildText("dd.bookmarks")
		Fanfic.Hits = e.ChildText("dd.hits")

	})
	c.OnHTML("#workskin > div.preface.group", func(e *colly.HTMLElement) {
		// log.Println(e.ChildText("h2.title.heading"))
		Fanfic.Title = e.ChildText("h2.title.heading")
		Fanfic.Author = e.ChildText("h3 > a")
		// Fanfic.Title = tmp

		// Title = e.ChildText("h2.title.heading")
		// Author = ""
	})

	c.OnHTML("div.summary.module", func(e *colly.HTMLElement) {
		//log.Println("Summary debug")
		//log.Println(len(e.Text))
		//log.Println(e.Text)
		var sum []string
		//var Summary string

		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			//log.Println(reflect.TypeOf(el.Text))
			txt := fmt.Sprintf("%s", el.Text)
			sum = append(sum, txt)
			//Stats.Summary = el.Text
		})
		Fanfic.Summary = sum

	})

	c.OnHTML("dd.fandom.tags", func(e *colly.HTMLElement) {
		fandom := e.ChildText("a.tag")
		if fandom == "" {
			log.Printf("Fandom is null")
		} //else {
		//fmt.Println(Fandom)
		//}
		Fanfic.Fandom = fandom
	})
	c.OnHTML("dd.relationship.tags", func(e *colly.HTMLElement) {
		var relationships []string
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			relationship := el.Text
			relationships = append(relationships, relationship)
		})
		//relationship := strings.Join(relationships, " | ")
		Fanfic.Relationship = relationships

	})

	c.OnHTML("dd.freeform.tags", func(e *colly.HTMLElement) {
		var AlternativeTags []string
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			if len(el.Text) != 0 {
				AlternativeTag := el.Text
				AlternativeTags = append(AlternativeTags, AlternativeTag)
			}
		})
		//AlternativeTag := strings.Join(AlternativeTags, " | ")
		//log.Println(AlternativeTag)
		Fanfic.AlternativeTags = AlternativeTags

	})

	c.Visit(Fanfic.URL)
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		// StatusCode := r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		// StatusCode := r.StatusCode
	})

	c.Wait()
	return Fanfic
}
