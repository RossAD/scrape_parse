package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"
	// "strconv"
	// "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
	// "io/ioutil"
)

type ResultObj struct {
	id          string
	repost_id   string
	title       string
	url         string
	price       string
	housing     string
	hood        string
	posted_time time.Time
}

func main() {
	postal := flag.Int("postal", 94118, "Zip Code of search area")
	distance := flag.Int("distance", 3, "Search radius away from zip code")
	min_price := flag.Int("min_price", 1000, "Min price to search for")
	max_price := flag.Int("max_price", 2500, "Max price to search for")

	database, _ := sql.Open("sqlite3", "./db/craigslist.db")
	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS results 
					(id INTEGER PRIMARY KEY, repost_id INTEGER, 
					title TEXT, url TEXT, price TEXT, housing TEXT, 
					hood TEXT, posted_time DATETIME,
					mod_time DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	statement.Exec()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("craigslist.org", "sfbay.craigslist.org"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("li.result-row", func(e *colly.HTMLElement) {
		fmt.Printf("Time Stamp: ")
		posted_time := e.ChildAttr("time.result-date", "datetime")
		fmt.Printf(posted_time + "\n")
		shortForm := "2006-01-02 15:04"
		loc, _ := time.LoadLocation("America/Los_Angeles")
		post_time, _ := time.ParseInLocation(shortForm, posted_time, loc)
		fmt.Printf("Posted Time: ")
		fmt.Println(post_time.In(loc))
		result_obj := ResultObj{
			id:          e.Attr("data-pid"),
			repost_id:   e.Attr("data-repost-of"),
			title:       e.ChildText("a.result-title"),
			url:         e.ChildAttr("a.result-title", "href"),
			posted_time: post_time,
		}
		e.ForEach("span.result-meta", func(_ int, ei *colly.HTMLElement) {
			price := ei.ChildText("span.result-price")
			housing := ei.ChildText("span.housing")
			hood := ei.ChildText("span.result-hood")
			result_obj.price = price
			result_obj.housing = housing
			result_obj.hood = hood
		})
		row := database.QueryRow(fmt.Sprintf("SELECT id FROM results where id='%s'", result_obj.id))
		var id string
		row.Scan(&id)
		if (id != result_obj.id) && (id != result_obj.repost_id) {
			fmt.Printf("New Record Insert!!!!!!!!\n")
			statement, _ := database.Prepare(`INSERT INTO results (id, repost_id, title, url, price, posted_time, housing, hood)
			   VALUES (?, ?, ?, ? ,?, ?, ?, ?)`)
			statement.Exec(result_obj.id, result_obj.repost_id, result_obj.title, result_obj.url,
				result_obj.price, result_obj.posted_time, result_obj.housing, result_obj.hood)
		}
		rows, _ := database.Query("SELECT id, title, price FROM results")
		var title string
		var price string
		for rows.Next() {
			rows.Scan(&id, &title, &price)
			fmt.Println(id + " : " + title + "\nPrice: " + price)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://sfbay.craigslist.org/search/nby/apa?search_distance=5&postal=94941&max_price=2500&availabilityMode=0&sale_date=all+date")

}
