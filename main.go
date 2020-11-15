package ao3api

import (
	"fmt"
	"strings"
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

// ParseWork Parse ao3 work from wID and cID
func ParseWork(wID, cID string) (cTitle, Title, Author string, chaps []string) {
	var sWork work
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", wID)
	title := ""
	author := ""
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	c.Limit(&colly.LimitRule{ //nolint
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
		//fmt.Println(chaps[0])
		//fmt.Println(FindChapters(cID, cIDs, chapsText))
		sWork.Author = author
		sWork.Title = title
		sWork.wID = wID
		sWork.cID = cID
		sWork.cTitle = cTitle
		sWork.Chaps = chaps
		//fmt.Printf("Title is %s Chapters title is %s\n", title, cTitle)
	})
	c.OnHTML("", func(e *colly.HTMLElement) {
		//fmt.Println(e.Text)
	})

	//fmt.Printf("Title is %s Chapters title is %s\n", title, cTitle)
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})
	for i := 0; i < 5; i++ {
		c.Visit(url) //nolint
	}
	c.Wait()
	//fmt.Println(sWork.cTitle, title, author, sWork.Chaps)
	return sWork.cTitle, sWork.Title, sWork.Author, sWork.Chaps
}

//ParseSummary to find summary of work
func ParseSummary(wID, cID string) string {
	var Summary string
	cTitle, title, author, chaps := ParseWork(wID, cID)
	_, _, _ = cTitle, title, author
	url := fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s", wID, chaps[0])
	//fmt.Println(url)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	c.Limit(&colly.LimitRule{ //nolint
		// Filter domains affected by this rule
		DomainGlob: "*archiveofourown.org/*",
		// Set a delay between requests to these domains
		Delay: 5 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
		// Add User Agent
		Parallelism: 2,
	})

	c.OnHTML("blockquote", func(e *colly.HTMLElement) {
		Summary = strings.Trim(e.Text, "\t \n")
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})
	for i := 0; i < 5; i++ {
		c.Visit(url) //nolint
	}
	c.Wait()

	return Summary
}
