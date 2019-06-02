package main

import (
	// "database/sql"
	"fmt"
	// "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	// _ "github.com/mattn/go-sqlite3"
	// "io/ioutil"
)

type ResultObj struct {
	id        string
	repost_id string
	title     string
	url       string
	price     string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("craigslist.org", "sfbay.craigslist.org"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("li.result-row", func(e *colly.HTMLElement) {
		result_obj := ResultObj{
			id:        e.Attr("data-pid"),
			repost_id: e.Attr("data-repost-of"),
			title:     e.ChildText("a.result-title"),
			url:       e.ChildAttr("a.result-title", "href"),
		}
		e.ForEach("a.result-image", func(_ int, ei *colly.HTMLElement) {
			price := ei.ChildText("span.result-price")
			result_obj.price = price
		})
		fmt.Printf("Post Id Found: %s\n", result_obj.id)
		fmt.Printf("Post Repost Id Found: %s\n", result_obj.repost_id)
		fmt.Printf("Post Title Found: %s\n", result_obj.title)
		fmt.Printf("Post url: %s\n", result_obj.url)
		fmt.Printf("Post price: %s\n", result_obj.price)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://sfbay.craigslist.org/search/apa?search_distance=3&postal=94118&max_price=2000&availabilityMode=0&sale_date=all+dates")

}
