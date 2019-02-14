package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("https://sfbay.craigslist.org/search/apa?search_distance=3&postal=94118&max_price=2000&min_bathrooms=1&availabilityMode=0&sale_date=all+dates")
	if err != nil {
		fmt.Printf("Error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("Error2")
		}
		bodyString := string(bodyBytes)
		fmt.Printf(bodyString)
	}

}
