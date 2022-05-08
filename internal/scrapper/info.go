package scrapper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/capoverflow/ao3api/models"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"

	"github.com/corpix/uarand"
)

func GetInfo(Params models.FanficParams, ChaptersIDs []string) models.Work {
	var Fanfic models.Work
	Fanfic.ChaptersIDs = ChaptersIDs
	Fanfic.URL = fmt.Sprintf("http://%s/works/%s/chapters/%s?view_adult=true", Params.Addr, Params.WorkID, Fanfic.ChaptersIDs[0])
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.UserAgent(uarand.GetRandom()),
		//colly.AllowURLRevisit(),
	)
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		// DomainGlob: "*archiveofourown.org/*",
		// Set a delay between requests to these domains
		Delay: 15 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
		// Add User Agent
		Parallelism: 2,
	})

	if len(Params.ProxyURLs) != 0 {
		log.Println("using proxy")
		rp, err := proxy.RoundRobinProxySwitcher(Params.ProxyURLs...)
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

	// github.com/capoverflow/ao3api: 2022/04/23 02:45:35 info.go:102: https://download.archiveofourown.org/downloads/33638854/Trial%20and%20Error.azw3?updated_at=1650665615 <nil>
	// github.com/capoverflow/ao3api: 2022/04/23 02:45:35 info.go:102: https://download.archiveofourown.org/downloads/33638854/Trial%20and%20Error.epub?updated_at=1650665615 <nil>
	// github.com/capoverflow/ao3api: 2022/04/23 02:45:35 info.go:102: https://download.archiveofourown.org/downloads/33638854/Trial%20and%20Error.mobi?updated_at=1650665615 <nil>
	// github.com/capoverflow/ao3api: 2022/04/23 02:45:35 info.go:102: https://download.archiveofourown.org/downloads/33638854/Trial%20and%20Error.pdf?updated_at=1650665615 <nil>
	// github.com/capoverflow/ao3api: 2022/04/23 02:45:35 info.go:102: https://download.archiveofourown.org/downloads/33638854/Trial%20and%20Error.html?updated_at=1650665615 <nil>

	c.OnHTML("li.download", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			if !strings.Contains(el.Attr("href"), "#") {
				href := fmt.Sprintf("https://download.archiveofourown.org%s", el.Attr("href"))

				FileType := []string{"azw3", "epub", "mobi", "pdf", "html"}
				for _, ft := range FileType {
					switch {
					case strings.Contains(href, ft):
						download := models.Downloads{
							FileType: ft,
							Url:      href,
						}
						Fanfic.Downloads = append(Fanfic.Downloads, download)
					}
				}

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
