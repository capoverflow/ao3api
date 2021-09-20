package scrapper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.com/capoverflow/ao3api/models"

	"github.com/gocolly/colly"

	"github.com/corpix/uarand"
)

func GetInfo(WorkID string, ChaptersIDs []string) models.Work {
	var Fanfic models.Work
	Fanfic.ChaptersIDs = ChaptersIDs
	Fanfic.URL = fmt.Sprintf("https://archiveofourown.org/works/%s/chapters/%s?view_adult=true", WorkID, Fanfic.ChaptersIDs[0])
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.UserAgent(uarand.GetRandom()),
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

		Fanfic = models.Work{
			Published: e.ChildText("dd.published"),
			Updated:   e.ChildText("dd.status"),
			Words:     e.ChildText("dd.words"),
			Chapters:  e.ChildText("dd.chapters"),
			Comments:  e.ChildText("dd.comments"),
			Kudos:     e.ChildText("dd.kudos"),
			Bookmarks: e.ChildText("dd.bookmarks"),
			Hits:      e.ChildText("dd.hits"),
		}

	})
	c.OnHTML("div.preface.group", func(e *colly.HTMLElement) {
		Fanfic.Title = e.ChildText("h2.title.heading")
		// Fanfic.Author = e.ChildText("h3.byline.heading")
		// log.Println(e.ChildAttrs("a", "href"))
	})

	c.OnHTML("h3.byline.heading", func(e *colly.HTMLElement) {
		// log.Println(e.ChildAttrs("a", "href"))
		// log.Println(e.ChildText("a"))
		// Fanfic.Author = e.ChildText("a")
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			// log.Println(h.Text)
			Fanfic.Author = append(Fanfic.Author, h.Text)
		})
	})

	c.OnHTML("div.summary.module", func(e *colly.HTMLElement) {

		var sum []string

		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			// txt := fmt.Sprintf("%s", el.Text)
			txt := el.Text
			sum = append(sum, txt)
		})
		Fanfic.Summary = sum

	})

	c.OnHTML("dd.fandom.tags", func(e *colly.HTMLElement) {
		var fandoms []string
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			fandom := el.Text
			fandoms = append(fandoms, fandom)

		})
		// log.Println(fandoms)
		Fanfic.Fandom = fandoms
	})
	c.OnHTML("dd.relationship.tags", func(e *colly.HTMLElement) {
		var relationships []string
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			relationship := el.Text
			relationships = append(relationships, relationship)
		})
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
		Fanfic.AlternativeTags = AlternativeTags

	})

	c.OnHTML("li.download", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			// log.Println(el.Attr("href"))
			if !strings.Contains(el.Attr("href"), "#") {
				Fanfic.Downloads = append(Fanfic.Downloads, fmt.Sprintf("https://download.archiveofourown.org%s", el.Attr("href")))
			}
			// log.Println(a)
			// log.Printf("https://download.archiveofourown.org%s", el.Attr("href"))
		})
	})

	c.Visit(Fanfic.URL)
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.Wait()
	return Fanfic
}
