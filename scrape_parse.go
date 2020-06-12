package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	// "strconv"
	// "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
	// "io/ioutil"
)

// ResultObj is a data structure for storing the results of the page scrape
type ResultObj struct {
	id         string
	repostId   string
	title      string
	url        string
	price      string
	housing    string
	hood       string
	postedTime time.Time
}

func dbWrite(listing ResultObj) {
	database, err := sql.Open("sqlite3", "./db/craigslist.db")
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	defer database.Close()
	statement, err := database.Prepare(`CREATE TABLE IF NOT EXISTS results 
					(id INTEGER PRIMARY KEY, repostId INTEGER, 
					title TEXT, url TEXT, price TEXT, housing TEXT, 
					hood TEXT, postedTime DATETIME,
					mod_time DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatal("Failed to create table ", err)
	}

	statement.Exec()
	row := database.QueryRow(fmt.Sprintf("SELECT id FROM results where id='%s'", listing.id))
	var id string
	row.Scan(&id)
	if (id != listing.id) && (id != listing.repostId) {
		fmt.Printf("New Record Insert!!!!!!!!\n")
		statement, err := database.Prepare(`INSERT INTO results (id, repostId, title, url, price, postedTime, housing, hood)
			   VALUES (?, ?, ?, ? ,?, ?, ?, ?)`)
		if err != nil {
			log.Fatal("Failed to create insert statement ", err)
		}
		result, err := statement.Exec(listing.id, listing.repostId, listing.title, listing.url,
			listing.price, listing.postedTime, listing.housing, listing.hood)
		if err != nil {
			log.Fatal("Failed to insert values ", err)
		}
		var count int
		count = 0
		if result != nil {
			count++
		}

	}
	rows, _ := database.Query("SELECT id, title, price FROM results")
	var title string
	var price string
	for rows.Next() {
		rows.Scan(&id, &title, &price)
	}
}

func main() {
	postal := flag.Int("postal", 94941, "Zip Code of search area")
	distance := flag.Int("distance", 2, "Search radius away from zip code")
	// min_price := flag.Int("min_price", 1000, "Min price to search for")
	maxPrice := flag.Int("max_price", 2500, "Max price to search for")

	flag.Parse()
	if postal == nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	urlBase := "https://sfbay.craigslist.org/search/nby/apa"
	fmt.Printf("%d %d %d\n", *postal, *distance, *maxPrice)
	queryString := fmt.Sprintf("?search_distance=%d&postal=%d&max_price=%d&availabilityMode=0&sale_date=all+date", *distance, *postal, *maxPrice)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("craigslist.org", "sfbay.craigslist.org"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("li.result-row", func(e *colly.HTMLElement) {
		// fmt.Printf("Time Stamp: ")
		postedTime := e.ChildAttr("time.result-date", "datetime")
		// fmt.Printf(postedTime + "\n")
		shortForm := "2006-01-02 15:04"
		loc, _ := time.LoadLocation("America/Los_Angeles")
		postTime, _ := time.ParseInLocation(shortForm, postedTime, loc)
		// fmt.Printf("Posted Time: ")
		// fmt.Println(postTime.In(loc))
		resultObj := ResultObj{
			id:         e.Attr("data-pid"),
			repostId:   e.Attr("data-repost-of"),
			title:      e.ChildText("a.result-title"),
			url:        e.ChildAttr("a.result-title", "href"),
			postedTime: postTime,
		}
		e.ForEach("span.result-meta", func(_ int, ei *colly.HTMLElement) {
			price := ei.ChildText("span.result-price")
			housing := ei.ChildText("span.housing")
			hood := ei.ChildText("span.result-hood")
			resultObj.price = price
			resultObj.housing = housing
			resultObj.hood = hood
		})
		dbWrite(resultObj)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(urlBase + queryString)

}
