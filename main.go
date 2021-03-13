package ao3

//package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

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

// Works ...
func Works(wID, cID string) Work {
	var fanfic Work
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", wID)
	//url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate", wID)
	// var title string
	// var author string
	//var cIDs []string
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

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		// title := e.ChildText("a:nth-child(1)")
		// log.Println(reflect.TypeOf(title))
		fanfic.Title = e.ChildText("a:nth-child(1)")
		fanfic.Author = e.ChildText("a:nth-child(2)")
	})

	c.OnHTML("#main > ol", func(e *colly.HTMLElement) {
		cIDs := e.ChildAttrs("a", "href")

		chapsText := []string{}
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			chapsText = append(chapsText, el.Text)
			//log.Println(el.Text)
		})

		cTitle, ChapterIDs, chapsText := FindChapters(cID, cIDs, chapsText)
		// fanfic.Author = author
		// fanfic.Title = title
		fanfic.WorkID = wID
		fanfic.ChapterID = cID
		fanfic.ChapterTitle = cTitle
		fanfic.ChaptersIDs = ChapterIDs
		fanfic.ChaptersTitles = chapsText
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Fatal(r.Request)
			log.Println(string(r.Body))
		}
	})
	//for i := 0; i < 4; i++ {
	//	c.Visit(url)
	//}
	c.Visit(url)
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	//log.Println(fanfic)
	return fanfic
}

//Info ...
func Info(wID, cID string) Work {
	var Stats Work
	//log.Println("Info", wID, cID)
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s?view_adult=true", wID, cID)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	//log.Println("url", url)
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
		Stats.Published = e.ChildText("dd.published")
		Stats.Updated = e.ChildText("dd.status")
		Stats.Words = e.ChildText("dd.words")
		Stats.Chapters = e.ChildText("dd.chapters")
		Stats.Comments = e.ChildText("dd.comments")
		Stats.Kudos = e.ChildText("dd.kudos")
		Stats.Bookmarks = e.ChildText("dd.bookmarks")
		Stats.Hits = e.ChildText("dd.hits")

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
		Stats.Summary = sum

	})
	c.OnHTML("dd.fandom.tags", func(e *colly.HTMLElement) {
		fandom := e.ChildText("a.tag")
		if fandom == "" {
			log.Printf("Fandom is null")
		} //else {
		//fmt.Println(Fandom)
		//}
		Stats.Fandom = fandom
	})
	c.OnHTML("dd.relationship.tags", func(e *colly.HTMLElement) {
		var relationships []string
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			relationship := el.Text
			relationships = append(relationships, relationship)
		})
		//relationship := strings.Join(relationships, " | ")
		Stats.Relationship = relationships

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
		Stats.AlternativeTags = AlternativeTags

	})
	c.OnRequest(func(r *colly.Request) {
		//log.Println("visiting", r.URL.String())

	})

	c.Visit(url)
	c.Wait()
	return Stats
}
