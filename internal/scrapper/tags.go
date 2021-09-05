package scrapper

import (
	"log"
	"strings"
)

func Tags(tags string) string {

	log.Println(tags)

	log.Println(strings.Split(tags, "a"))

	return tags

	// // path := "path with?reserved+characters"
	// // log.Println(url.PathEscape(path))

	// url := fmt.Sprintf("https://archiveofourown.org/tags/%s/works", url.PathEscape(tags))
	// // log.Printf("WorkID: %s, url %s", WorkID, url)
	// c := colly.NewCollector(
	// 	colly.CacheDir("./cache"),
	// 	colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
	// 	colly.AllowURLRevisit(),
	// )
	// c.Limit(&colly.LimitRule{
	// 	// Filter domains affected by this rule
	// 	DomainGlob: "*archiveofourown.org/*",
	// 	// Set a delay between requests to these domains
	// 	Delay: 5 * time.Second,
	// 	// Add an additional random delay
	// 	RandomDelay: 10 * time.Second,
	// 	// Add User Agent
	// 	Parallelism: 2,
	// })

	// c.OnRequest(func(r *colly.Request) {

	// })
	// c.OnScraped(func(r *colly.Response) { // DONE
	// 	if len(r.Body) == 0 {
	// 		log.Println(r.Request)
	// 		log.Println(string(r.Body))
	// 	}
	// })

	// // extract status code
	// c.OnResponse(func(r *colly.Response) {
	// 	// log.Println("response received", r.StatusCode)
	// })
	// c.OnError(func(r *colly.Response, err error) {
	// 	// log.Println("error:", r.StatusCode, err)
	// })

	// c.Visit(url)

	// return url
}
