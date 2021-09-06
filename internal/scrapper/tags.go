package scrapper

import (
	"log"
	"net/url"
	"time"

	"github.com/gocolly/colly"
)

func Search(search string) {

	searchURL := "https://archiveofourown.org/works/search?utf8=✓"
	// https://archiveofourown.org/works/search?utf8=✓&work_search[query]=Clarke griffin
	u, err := url.Parse(searchURL)
	if err != nil {
		log.Println(err)
	}
	log.Println(u.Query())
	q := u.Query()
	q.Set("work_search[query]", "Clarke")
	log.Println(u.String())
	u.RawQuery = q.Encode()
	log.Println(u.String())
	// tags = u.Path
	// return search

	// // path := "path with?reserved+characters"
	// // log.Println(url.PathEscape(path))

	// url := fmt.Sprintf("https://archiveofourown.org/tags/%s/works", url.PathEscape(tags))
	// // log.Printf("WorkID: %s, url %s", WorkID, url)
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

	c.OnHTML("ol.work.index.group", func(e *colly.HTMLElement) {
		// log.Println(e.Text)
		// stuff = e.Text
		e.ForEach("li > div.header.module > h4.heading", func(_ int, el *colly.HTMLElement) {

			// el.ForEach("a", func(_ int, em *colly.HTMLElement) {
			// 	// link, _ := em.DOM.Find("a").Attr("href")
			// 	// // link = fmt.Sprintf("https://archiveofourown.org%s", link)
			// 	// // log.Println(em.DOM.Find("a").Text(), link)
			// 	// log.Println(link)
			// 	log.Println(em.Text)

			// })
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
