//package neogoapi
package main

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

func ParseWorks(wID, cID string) { //(ChaptersTitles, WorkTitle, WorkAuthor string, WorkChapters []string) {
	//var sWork work

	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", wID)
	var title string
	var author string
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

	c.Visit(url)
	c.Wait()
	fmt.Println(title, author)
}
func main() {

	ParseWorks("4854050", "")
}
