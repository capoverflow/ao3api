package fanfic

import (
	"ao3api/internals/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/corpix/uarand"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

// func
// GetInfo retrieves information about a fanfic work from a given URL
// using the provided Params and ChaptersIDs.
func GetInfo(Params models.FanficParams, ChaptersIDs []string) (Fanfic models.Fanfic) {
	// Initialize an empty Work struct
	// Set the ChaptersIDs field of the Work struct to the provided ChaptersIDs
	Fanfic.ChaptersIDs = ChaptersIDs
	// Set the URL field of the Work struct to the URL of the first chapter in the ChaptersIDs slice
	Fanfic.URL = fmt.Sprintf("http://%s/works/%s/chapters/%s?view_adult=true", Params.Addr, Params.WorkID, Fanfic.ChaptersIDs[0])

	// Initialize a new colly collector with the specified options
	c := colly.NewCollector(
		// Use the provided cache directory
		colly.CacheDir("./cache"),
		// Use a random user agent
		colly.UserAgent(uarand.GetRandom()),
	)
	// Set a rate limit for requests to the specified domains
	c.Limit(&colly.LimitRule{
		// Affected domains
		DomainGlob: "*archiveofourown.org/*",
		// Delay between requests
		Delay: 15 * time.Second,
		// Additional random delay
		RandomDelay: 10 * time.Second,
		// Set the maximum number of concurrent requests to 2
		Parallelism: 2,
	})

	// If there are proxy URLs provided, use them for requests
	if len(Params.ProxyURLs) != 0 {
		log.Println("using proxy")
		// Create a proxy switcher that uses a round-robin algorithm to choose the next proxy
		rp, err := proxy.RoundRobinProxySwitcher(Params.ProxyURLs...)
		if err != nil {
			log.Fatal(err)
		}
		// Set the proxy switcher as the proxy function for the collector
		c.SetProxyFunc(rp)
	}

	// Set a handler for HTML elements matching the "dl.stats" selector
	c.OnHTML("dl.stats", func(e *colly.HTMLElement) {
		// Set the fields of the Work struct based on the values of the child elements of the "dl.stats" element
		Fanfic = models.Fanfic{
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
	// Set a handler for HTML elements matching the "h2.title.heading" selector
	c.OnHTML("h2.title.heading", func(e *colly.HTMLElement) {
		// Set the Title field of the Work struct to the text of the element, with leading and trailing whitespace trimmed
		Fanfic.Title = strings.TrimSpace(e.Text)
	})

	// Set a handler for HTML elements matching the "h3.byline.heading" selector
	c.OnHTML("h3.byline.heading", func(e *colly.HTMLElement) {
		// For each "a" element within the "h3.byline.heading" element,
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			// append the text of the element to the Author field of the Work struct
			Fanfic.Author = append(Fanfic.Author, h.Text)
		})
	})

	// Set a handler for HTML elements matching the "div.summary.module" selector
	c.OnHTML("div.summary.module", func(e *colly.HTMLElement) {
		// For each "p" element within the "div.summary.module" element,
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			// append the text of the element to the Summary field of the Work struct
			Fanfic.Summary = append(Fanfic.Summary, el.Text)
		})
	})

	// Set a handler for HTML elements matching the "dd.fandom.tags" selector
	c.OnHTML("dd.fandom.tags", func(e *colly.HTMLElement) {
		// For each "a.tag" element within the "dd.fandom.tags" element,
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			// append the text of the element to the Fandom field of the Work struct
			Fanfic.Fandom = append(Fanfic.Fandom, el.Text)
		})
	})
	// Set a handler for HTML elements matching the "dd.relationship.tags" selector
	c.OnHTML("dd.relationship.tags", func(e *colly.HTMLElement) {
		// For each "a.tag" element within the "dd.relationship.tags" element,
		e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
			// append the text of the element to the Relationship field of the Work struct
			Fanfic.Relationship = append(Fanfic.Relationship, el.Text)
		})
	})

	// // Set a handler for HTML elements matching the "dd.freeform.tags" selector
	// c.OnHTML("dd.freeform.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.freeform.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Freeform field of the Work struct
	// 		Fanfic.Freeform = append(Fanfic.Freeform, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.rating.tags" selector
	// c.OnHTML("dd.rating.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.rating.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Rating field of the Work struct
	// 		Fanfic.Rating = append(Fanfic.Rating, el.Text)
	// 	})
	// })
	// // Set a handler for HTML elements matching the "dd.character.tags" selector
	// c.OnHTML("dd.character.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.character.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Character field of the Work struct
	// 		Fanfic.Character = append(Fanfic.Character, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.warnings.tags" selector
	// c.OnHTML("dd.warnings.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.warnings.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Warning field of the Work struct
	// 		Fanfic.Warning = append(Fanfic.Warning, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.category.tags" selector
	// c.OnHTML("dd.category.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.category.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Category field of the Work struct
	// 		Fanfic.Category = append(Fanfic.Category, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.language.tags" selector
	// c.OnHTML("dd.language.tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.language.tags" element,
	// 	// For each "a.tag" element within the "dd.language.tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Language field of the Work struct
	// 		Fanfic.Language = append(Fanfic.Language, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.additional_tags" selector
	// c.OnHTML("dd.additional_tags", func(e *colly.HTMLElement) {
	// 	// For each "a.tag" element within the "dd.additional_tags" element,
	// 	e.ForEach("a.tag", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Additional field of the Work struct
	// 		Fanfic.Additional = append(Fanfic.Additional, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.collections" selector
	// c.OnHTML("dd.collections", func(e *colly.HTMLElement) {
	// 	// For each "a" element within the "dd.collections" element,
	// 	e.ForEach("a", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Collection field of the Work struct
	// 		Fanfic.Collection = append(Fanfic.Collection, el.Text)
	// 	})
	// })

	// // Set a handler for HTML elements matching the "dd.challenges" selector
	// c.OnHTML("dd.challenges", func(e *colly.HTMLElement) {
	// 	// For each "a" element within the "dd.challenges" element,
	// 	e.ForEach("a", func(_ int, el *colly.HTMLElement) {
	// 		// append the text of the element to the Challenge field of the Work struct
	// 		Fanfic.Challenge = append(Fanfic.Challenge, el.Text)
	// 	})
	// })

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

	// Visit the URL of the first chapter in the ChaptersIDs slice to retrieve the information
	c.Visit(Fanfic.URL)
	// Return the Work struct
	return Fanfic
}
