package ao3

//package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type work struct {
	wID        string   //`json:"wID,omitempty"`
	cID        string   //`json:"cID,omitempty"`
	cTitle     string   //`json:"cTitle,omitempty"`
	Title      string   //`json:"Title,omitempty"`
	Author     string   //`json:"Author,omitempty"`
	Chaps      []string //`json:"ChaptersTitles,omitempty"`
	ChapterIDs []string //
}

type id struct {
	WorkID    string
	ChapterID string
}

type ids struct {
	works []id
}

type stats struct {
	Published       string
	Updated         string
	Words           string
	Chapters        string
	Comments        string
	Kudos           string
	Bookmarks       string
	Hits            string
	Summary         string
	Fandom          string
	Relationship    []string
	AlternativeTags []string
}

// Works ...
func Works(wID, cID string) (ChaptersTitles, WorkTitle, WorkAuthor string, ChapterIDs []string, Chaps []string) {
	var sWork work
	//url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", wID)
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate", wID)
	var title string
	var author string
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
		title = e.ChildText("a:nth-child(1)")
		author = e.ChildText("a:nth-child(2)")
	})

	c.OnHTML("#main > ol", func(e *colly.HTMLElement) {
		cIDs := e.ChildAttrs("a", "href")
		chapsText := []string{}
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			chapsText = append(chapsText, el.Text)
			//log.Println(el.Text)
		})

		cTitle, ChapterIDs, chapsText := FindChapters(cID, cIDs, chapsText)
		//log.Println(chaps[0])
		//log.Println(FindChapters(cID, cIDs, chapsText))
		sWork.Author = author
		sWork.Title = title
		sWork.wID = wID
		sWork.cID = cID
		sWork.cTitle = cTitle
		sWork.ChapterIDs = ChapterIDs
		sWork.Chaps = chapsText
		//log.Println(chaps)
		//fmt.Printf("Title is %s Chapters title is %s\n", title, cTitle)
	})
	//c.OnHTML("*", func(e *colly.HTMLElement) {
	//	ao3Error = e
	//	fmt.Println([]byte(ao3Error))
	//})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Fatal(r.Request)
			log.Println(string(r.Body))
		}
		//fmt.Println(len(string(r.Body)))

	})
	//for i := 0; i < 4; i++ {
	//	c.Visit(url)
	//}
	c.Visit(url)
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//log.Println(sWork.Chaps)
	return sWork.cTitle, sWork.Title, sWork.Author, sWork.ChapterIDs, sWork.Chaps
}

//Info ...
func Info(wID, cID string) (Published, Updated, Words, Chapters, Comments, Kudos, Bookmarks, Hits, Summary, Fandom string, Relationship, AlternativeTags []string) {
	var Stats stats
	//log.Println(wID, cID)
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s?view_adult=true", wID, cID)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
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
		//log.Println(len(e.Text))
		var sum []string
		var Summary string

		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			//fmt.Println(el.Text)
			sum = append(sum, el.Text)
			//Stats.Summary = el.Text
		})
		//if len(sum) == 1 {
		//	Summary = sum[0]
		//} else if len(sum) == 2 {
		//	Summary = fmt.Sprintf("%s %s", sum[0], sum[1])
		//} else {
		//	log.Println("Error in summary")
		//}

		Summary = strings.Join(sum, " ")

		//Summary = fmt.Sprintf("%q\n", sum) //
		//Summary = fmt.Sprintf("%s %s", sum[0], sum[1])
		//log.Println(len(sum))
		//log.Println(sum)
		Stats.Summary = Summary

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
		//fmt.Println(AlternativeTag)
		Stats.AlternativeTags = AlternativeTags

	})
	c.OnRequest(func(r *colly.Request) {
		//log.Println("visiting", r.URL.String())

	})

	c.Visit(url)
	c.Wait()
	return Stats.Published, Stats.Updated, Stats.Words, Stats.Chapters, Stats.Comments, Stats.Kudos, Stats.Bookmarks, Stats.Hits, Stats.Summary, Stats.Fandom, Stats.Relationship, Stats.AlternativeTags
}
