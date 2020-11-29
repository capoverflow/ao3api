package ao3

//package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type work struct {
	wID    string   //`json:"wID,omitempty"`
	cID    string   //`json:"cID,omitempty"`
	cTitle string   //`json:"cTitle,omitempty"`
	Title  string   //`json:"Title,omitempty"`
	Author string   //`json:"Author,omitempty"`
	Chaps  []string //`json:"ChaptersTitles,omitempty"`
}

type id struct {
	WorkID    string
	ChapterID string
}

type ids struct {
	works []id
}

func ParseWorks(wID, cID string) (ChaptersTitles, WorkTitle, WorkAuthor string, WorkChapters []string) {
	var sWork work

	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", wID)
	var title string
	var author string
	//var cIDs []string
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

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		title = e.ChildText("a:nth-child(1)")
		author = e.ChildText("a:nth-child(2)")
	})

	c.OnHTML("#main > ol", func(e *colly.HTMLElement) {
		cIDs := e.ChildAttrs("a", "href")
		chapsText := []string{}
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			chapsText = append(chapsText, el.Text)
			//fmt.Println(el.Text)
		})

		cTitle, chaps := FindChapters(cID, cIDs, chapsText)
		fmt.Println(chaps[0])
		//fmt.Println(FindChapters(cID, cIDs, chapsText))
		sWork.Author = author
		sWork.Title = title
		sWork.wID = wID
		sWork.cID = cID
		sWork.cTitle = cTitle
		sWork.Chaps = chaps
		//fmt.Println(chaps)
		//fmt.Printf("Title is %s Chapters title is %s\n", title, cTitle)
	})
	c.Visit(url)
	c.Wait()
	//fmt.Println(sWork.Chaps)
	return sWork.cTitle, sWork.Title, sWork.Author, sWork.Chaps
}

//ParseSummary
func ParseSummary(wID, cID string) {
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s", wID, cID)
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
	c.OnHTML("div", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.Visit(url)
	c.Wait()
}
