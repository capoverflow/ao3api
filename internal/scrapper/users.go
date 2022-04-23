package scrapper

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"gitlab.com/capoverflow/ao3api/internal/utils"
	"gitlab.com/capoverflow/ao3api/models"
)

func GetUsersInfo(User models.UserParams) (AuthorInfo models.User) {
	url := fmt.Sprintf("http://%s/users/%s/profile", User.Addr, User.Username)
	// log.Printf("WorkID: %s, url %s", WorkID, url)
	c := colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.UserAgent(uarand.GetRandom()),
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"),
		colly.AllowURLRevisit(),
	)
	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		// DomainGlob: "*archiveofourown.org/*",
		// Set a delay between requests to these domains
		Delay: 5 * time.Second,
		// Add an additional random delay
		RandomDelay: 10 * time.Second,
		// Add User Agent
		Parallelism: 2,
	})

	c.OnHTML("dl.meta", func(e *colly.HTMLElement) {

		AuthorInfo.Profile.Pseuds = e.ChildText("dd.pseuds")
		AuthorInfo.Profile.JoinDate = e.ChildText("dd:nth-child(4)")
		AuthorInfo.Profile.Email = e.ChildText("dd.email")
	})

	c.OnHTML("div.bio.module", func(e *colly.HTMLElement) {
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

func GetUserBookmarksPage(Author string) (PageNB int) {
	u, err := url.Parse(fmt.Sprintf("https://archiveofourown.org/users/%s/bookmarks", Author))

	if err != nil {
		log.Println(err)
	}

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

	c.OnHTML("ol.pagination.actions", func(e *colly.HTMLElement) {
		var links []string
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			// links = append(links, el.ChildAttr("a", "href"))
			if len(el.ChildAttr("a", "href")) != 0 {
				links = append(links, el.ChildAttr("a", "href"))

			}
		})
		PageNB = utils.FindUrl(links)

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

	c.Visit(u.String())

	return PageNB
}

func GetUserBookmarks(Author string) (Bookmarks []string) {
	u, err := url.Parse(fmt.Sprintf("https://archiveofourown.org/users/%s/bookmarks", Author))

	if err != nil {
		log.Println(err)
	}

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

	PageNB := GetUserBookmarksPage(Author)
	log.Println(PageNB)

	if PageNB == 0 {

		c.OnHTML("ol.bookmark.index.group", func(e *colly.HTMLElement) {
			e.ForEach("li.bookmark.blurb.group", func(_ int, el *colly.HTMLElement) {
				html, err := el.DOM.Html()
				if err != nil {
					log.Println(err)
				}
				dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
				if err != nil {
					log.Println(err)
				}
				html, err = dom.Find("h4.heading").Html()
				if err != nil {
					log.Println(err)
				}
				dom, err = goquery.NewDocumentFromReader(strings.NewReader(html))
				if err != nil {
					log.Println(err)
				}
				link, exist := dom.Find("a").Attr("href")
				if exist {
					Bookmarks = append(Bookmarks, link)
					// fmt.Printf("%s\n \n", link)
				}

			})

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
			// log.Println(r.Headers)

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

		c.Visit(u.String())

	} else {

		q, _ := queue.New(
			2, // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
		)

		for i := 1; i <= PageNB; i++ {
			if i == 1 {
				q.AddURL(u.String())
			} else {
				q.AddURL(fmt.Sprintf("%s?page=%d", u.String(), i))
			}

		}

		// log.Println(q.Size())

		c.OnHTML("ol.bookmark.index.group", func(e *colly.HTMLElement) {
			e.ForEach("li.bookmark.blurb.group", func(i int, el *colly.HTMLElement) {
				links := el.ChildAttrs("a", "href")
				re := regexp.MustCompile(`/works/[0-9]+/`)
				for _, link := range links {
					// fmt.Println(link)
					ff := re.FindString(link)
					if len(ff) != 0 {
						SplitString := strings.Split(ff, "/")
						Bookmarks = append(Bookmarks, SplitString[2])

						// log.Println(SplitString)
					}

				}

			})

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
			// log.Println(r.Headers)

		})

		q.Run(c)

	}

	return Bookmarks
}
