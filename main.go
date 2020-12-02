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

type stats struct {
	Published string
	Updated   string
	Words     string
	Chapters  string
	Comments  string
	Kudos     string
	Bookmarks string
	Hits      string
	Summary   string
}

func Works(wID, cID string) (ChaptersTitles, WorkTitle, WorkAuthor string, WorkChapters []string) {
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
			//log.Println(el.Text)
		})

		cTitle, chaps := FindChapters(cID, cIDs, chapsText)
		//log.Println(chaps[0])
		//log.Println(FindChapters(cID, cIDs, chapsText))
		sWork.Author = author
		sWork.Title = title
		sWork.wID = wID
		sWork.cID = cID
		sWork.cTitle = cTitle
		sWork.Chaps = chaps
		//log.Println(chaps)
		//fmt.Printf("Title is %s Chapters title is %s\n", title, cTitle)
	})
	c.Visit(url)
	c.Wait()
	//log.Println(sWork.Chaps)
	return sWork.cTitle, sWork.Title, sWork.Author, sWork.Chaps
}

//Info
func Info(wID, cID string) (Published, Updated, Words, Chapters, Comments, Kudos, Bookmarks, Hits, Summary string) {
	var Stats stats
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
	c.OnHTML("dl.stats", func(e *colly.HTMLElement) {
		var StatsText []string
		e.ForEach("dd", func(_ int, el *colly.HTMLElement) {
			//log.Println(el.Text, el.Index)
			StatsText = append(StatsText, el.Text)
		})
		Stats.Published = StatsText[0]
		Stats.Updated = StatsText[1]
		Stats.Words = StatsText[2]
		Stats.Chapters = StatsText[3]
		Stats.Comments = StatsText[4]
		Stats.Kudos = StatsText[5]
		Stats.Bookmarks = StatsText[6]
		Stats.Hits = StatsText[7]
		//log.Println(Stats)

	})

	c.OnHTML("div.summary.module", func(e *colly.HTMLElement) {
		var sum []string
		var Summary string
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			sum = append(sum, el.Text)
			//Stats.Summary = el.Text
		})
		Summary = fmt.Sprintf("%s %s", sum[0], sum[1])
		//log.Println(len(Summary))
		//log.Println(Summary)
		Stats.Summary = Summary

	})

	c.Visit(url)
	c.Wait()
	return Stats.Published, Stats.Updated, Stats.Words, Stats.Chapters, Stats.Comments, Stats.Kudos, Stats.Bookmarks, Stats.Hits, Stats.Summary
}
