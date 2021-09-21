package scrapper

import (
	"fmt"
	"log"
	"time"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"gitlab.com/capoverflow/ao3api/models"
)

func GetUsersInfo(Author string) (AuthorInfo models.User) {
	url := fmt.Sprintf("https://archiveofourown.org/users/%s/profile", Author)
	// log.Printf("WorkID: %s, url %s", WorkID, url)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.UserAgent(uarand.GetRandom()),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
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
	// if len(proxyURLs) != 0 {
	// 	rp, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	c.SetProxyFunc(rp)
	// }

	c.OnHTML("dl.meta", func(e *colly.HTMLElement) {
		// log.Println("Users")
		// log.Println(e.DOM.Html())
		AuthorInfo.Profile.Pseuds = e.ChildText("dd.pseuds")
		AuthorInfo.Profile.JoinDate = e.ChildText("dd:nth-child(4)")
		AuthorInfo.Profile.Email = e.ChildText("dd.email")
		// log.Println(e.ChildText("dd.pseuds"))
	})

	c.OnHTML("div.bio.module", func(e *colly.HTMLElement) {
		// log.Println(strings.TrimSpace(e.Text))
		AuthorInfo.Profile.Bio = e.Text
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
		log.Println(r.Headers)

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Println(r.Request)
			log.Println(string(r.Body))
		}

	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		fmt.Println(r.Ctx.Get("url"))

	})
	c.OnError(func(r *colly.Response, OnError error) {
		log.Println("error:", r.StatusCode)

	})

	c.Visit(url)
	return AuthorInfo
}

func GetUsersWorks(Author string) (AuthorInfo models.User) {
	url := fmt.Sprintf("https://archiveofourown.org/users/%s/works", Author)
	// log.Printf("WorkID: %s, url %s", WorkID, url)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.UserAgent(uarand.GetRandom()),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
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
	// if len(proxyURLs) != 0 {
	// 	rp, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	c.SetProxyFunc(rp)
	// }

	c.OnHTML("dl.meta", func(e *colly.HTMLElement) {

	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
		log.Println(r.Headers)

	})
	c.OnScraped(func(r *colly.Response) { // DONE
		if len(r.Body) == 0 {
			log.Println(r.Request)
			log.Println(string(r.Body))
		}

	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		fmt.Println(r.Ctx.Get("url"))

	})
	c.OnError(func(r *colly.Response, OnError error) {
		log.Println("error:", r.StatusCode)

	})

	c.Visit(url)
	return AuthorInfo
}
