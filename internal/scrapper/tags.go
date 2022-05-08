package scrapper

import (
	"log"
	"net/url"
	"time"

	"github.com/capoverflow/ao3api/models"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
)

func Search(SearchString models.Search) {
	log.Println(SearchString)

	searchURL := "https://archiveofourown.org/works/search?"
	u, err := url.Parse(searchURL)
	if err != nil {
		log.Println(err)
	}
	log.Println(u.Query())
	q := u.Query()
	if len(SearchString.AnyField) != 0 {
		q.Add("work_search[query]", SearchString.AnyField)
	}
	if len(SearchString.Title) != 0 {
		q.Add("work_search[title]", SearchString.Title)
	}
	if len(SearchString.Author) != 0 {
		q.Add("work_search[creators]", SearchString.Author)
	}
	if len(SearchString.Fandoms) != 0 {
		q.Add("work_search[fandom_names]", SearchString.Fandoms)
	}
	if len(SearchString.Relationship) != 0 {
		// Relationship := strings.Replace(SearchString.Relationship, "/", "*s*", -1)
		// log.Println(Relationship)
		q.Add("work_search[relationship_names]", SearchString.Relationship)
	}
	// q.Set("work_search[query]", )
	// q.Set("work_search[fandom_names]", SearchString.Fandoms)
	u.RawQuery = q.Encode()
	log.Println(q.Encode())

	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.UserAgent(uarand.GetRandom()),

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

	c.OnHTML("ol.work.index.group", func(e *colly.HTMLElement) {
		// log.Println(e.Text)
		// stuff = e.Text
		e.ForEach("li > div.header.module > h4.heading", func(_ int, el *colly.HTMLElement) {

			el.ForEach("a", func(_ int, em *colly.HTMLElement) {
				// link, _ := em.DOM.Find("a").Attr("href")
				// link = fmt.Sprintf("https://archiveofourown.org%s", link)
				// log.Println(em.DOM.Find("a").Text(), link)
				// log.Println(link)
				// log.Println(em.Text)

			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Do request", r.URL.String())

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Println(r.Request)
			log.Println(string(r.Body))
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		// log.Println("response received", r.StatusCode)
	})
	c.OnError(func(r *colly.Response, err error) {
		// log.Println("error:", r.StatusCode, err)
	})

	c.Visit(u.String())

	// return stuff
}
