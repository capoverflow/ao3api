package scrapper

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gitlab.com/capoverflow/ao3api/internal/utils"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

func GetFirstChapterID(WorkID, ChapterID string, proxyURLs []string, debug bool) (ChaptersIDs []string, StatusCode int, err error) {
	ChaptersIDs = []string{}
	err = nil

	url := fmt.Sprintf("https://archiveofourown.org/works/%s/navigate?view_adult=true", WorkID)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
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
	if len(proxyURLs) != 0 {
		rp, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
	if debug {
		c.OnHTML("html", func(e *colly.HTMLElement) {
			log.Println(e.Text)

		})
	} else {
		c.OnHTML("#signin", func(e *colly.HTMLElement) {
			err = errors.New("require login")
		})
		c.OnHTML("#main > ol", func(e *colly.HTMLElement) {
			hrefChaptersIDs := e.ChildAttrs("a", "href")
			ChaptersIDs = utils.FindChaptersIDs(hrefChaptersIDs)
		})
	}

	c.OnRequest(func(r *colly.Request) {
		if debug {
			log.Println("visiting", r.URL.String())
			log.Println(r.Headers)
		}
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
		StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, OnError error) {
		err = errors.New(OnError.Error())
		StatusCode = r.StatusCode
	})

	c.Visit(url)
	return ChaptersIDs, StatusCode, err
}
