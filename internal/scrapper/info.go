package scrapper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.com/capoverflow/ao3api/models"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"

	"github.com/corpix/uarand"
)

func GetInfo(WorkID string, ChaptersIDs []string, proxyURLs []string) models.Work {
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

	if len(proxyURLs) != 0 {
		rp, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}

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
	c.OnHTML("h2.title.heading", func(e *colly.HTMLElement) {
		Fanfic.Title = strings.TrimSpace(e.Text)
	})

	c.OnHTML("h3.byline.heading", func(e *colly.HTMLElement) {

		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			Fanfic.Author = append(Fanfic.Author, h.Text)
		})
	})

	c.OnHTML("div.summary.module", func(e *colly.HTMLElement) {

		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			Fanfic.Summary = append(Fanfic.Summary, el.Text)
		})

	})

	c.OnHTML("dd.fandom.tags", func(e *colly.HTMLElement) {
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			Fanfic.Fandom = append(Fanfic.Fandom, el.Text)
		})
	})
	c.OnHTML("dd.relationship.tags", func(e *colly.HTMLElement) {
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			Fanfic.Relationship = append(Fanfic.Relationship, el.Text)
		})
	})

	c.OnHTML("dd.freeform.tags", func(e *colly.HTMLElement) {
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			if len(el.Text) != 0 {
				Fanfic.AlternativeTags = append(Fanfic.AlternativeTags, el.Text)
			}
		})
	})

	c.OnHTML("li.download", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			if !strings.Contains(el.Attr("href"), "#") {
				Fanfic.Downloads = append(Fanfic.Downloads, fmt.Sprintf("https://download.archiveofourown.org%s", el.Attr("href")))
			}
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
