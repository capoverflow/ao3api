package scrapper

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/capoverflow/ao3api/internal/utils"
	"github.com/capoverflow/ao3api/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
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

	c.Visit(url)
	return AuthorInfo
}

func GetUsersWorks(User models.UserParams) (Works []string) {

	// https://archiveofourown.org/users/RhinoMouse/pseuds/Rhino/works

	// log.Println(User)
	u, err := url.Parse(fmt.Sprintf("https://%s/users/%s/pseuds/%s/works", User.Addr, User.Username, User.Pseuds))
	log.Println(u.String())
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

	PageNB := GetUsersWorksPage(User)
	log.Println(PageNB)

	if PageNB == 0 {

		c.OnHTML("ol.work.index.group", func(e *colly.HTMLElement) {
			e.ForEach("li.work.blurb.group", func(_ int, el *colly.HTMLElement) {
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
					Works = append(Works, link)
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
			// log.Println(r.Ctx.Get("url"))

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

		c.OnHTML("ol.work.index.group", func(e *colly.HTMLElement) {
			e.ForEach("li.work.blurb.group", func(i int, el *colly.HTMLElement) {
				links := el.ChildAttrs("a", "href")
				// log.Println(el)
				re := regexp.MustCompile(`/works/[0-9]+/`)
				for _, link := range links {

					// log.Println(link)
					ff := re.FindString(link)
					if len(ff) != 0 {
						SplitString := strings.Split(ff, "/")
						Works = append(Works, SplitString[2])

						// log.Println(SplitString)
					}

				}

			})

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
			// log.Println(r.Headers)

		})

		// c.OnResponse(func(r *colly.Response) {

		// 	log.Println("-----------------------------")

		// 	log.Println(r.StatusCode)

		// 	for key, value := range *r.Headers {
		// 		log.Printf("%s: %s\n", key, value)
		// 	}
		// })

		q.Run(c)

	}

	return Works
}

func GetUsersWorksPage(User models.UserParams) (PageNB int) {
	u, err := url.Parse(fmt.Sprintf("https://%s/users/%s/pseuds/%s/works", User.Addr, User.Username, User.Pseuds))
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

	// c.OnHTML("html", func(e *colly.HTMLElement) {
	// 	html, err := e.DOM.Html()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Printf("%s\n\n%s", User.Pseuds, html)
	// })

	c.OnHTML("ol.pagination.actions", func(e *colly.HTMLElement) {
		var links []string

		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			// links = append(links, el.ChildAttr("a", "href"))
			if len(el.ChildAttr("a", "href")) != 0 {
				links = append(links, el.ChildAttr("a", "href"))

			}
		})
		PageNB = utils.FindUrl("works", links)

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
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println("response received", r.StatusCode)
	// 	fmt.Println(r.Ctx.Get("url"))

	// })
	c.OnError(func(r *colly.Response, OnError error) {
		log.Println("error:", r.StatusCode)

	})

	c.Visit(u.String())

	return PageNB
}

func GetUserBookmarksPage(Params models.UserParams) (PageNB int) {
	u, err := url.Parse(fmt.Sprintf("https://%s/users/%s/bookmarks", Params.Addr, Params.Username))

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
		PageNB = utils.FindUrl("bookmarks", links)

	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
		// log.Println(r.Headers)
		// r.Headers.Set("cookie", "some cookie str")

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

func GetUserBookmarks(Params models.UserParams) (Bookmarks []string) {
	u, err := url.Parse(fmt.Sprintf("https://%s/users/%s/bookmarks", Params.Addr, Params.Username))

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

	PageNB := GetUserBookmarksPage(Params)
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
				// log.Println(el)
				re := regexp.MustCompile(`/works/[0-9]+/`)
				for _, link := range links {

					// log.Println(link)
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

		c.OnResponse(func(r *colly.Response) {

			log.Println("-----------------------------")

			log.Println(r.StatusCode)

			// for key, value := range *r.Headers {
			// 	log.Printf("%s: %s\n", key, value)
			// }
		})

		q.Run(c)

	}

	return Bookmarks
}
