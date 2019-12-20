package main

import (
	"fmt"

	"crypto/tls"
	"github.com/gocolly/colly"
	"net/http"

	utils "github.com/devopsext/utils"
)

var VERSION = "unknown"

var log = utils.GetLog()
var env = utils.GetEnvironment()

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 1,  so only  the  links on  the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(-1),
	)

	c.Async = false
	c.AllowedDomains = []string{"www.exness-168.com"}
	c.RedirectHandler = func(req *http.Request, via []*http.Request) error {

		log.Info("url: ", via[0].URL.String())
		fmt.Println("response: ", req.Response.StatusCode)
		return nil
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	c.OnResponse(func(response *colly.Response) {

		log.Info("url: ", response.Request.URL.String())
		fmt.Println("response: ", response.StatusCode)

	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println("visiting: ", e.Attr("href"))
		if err := e.Request.Visit(e.Attr("href")); err != nil {
			// Ignore already visited error, this appears too often
			if err != colly.ErrAlreadyVisited {
				//fmt.Println("error while linking: ", err.Error())
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	/*c.OnResponse(func(r *colly.Response) {
		fmt.Println("Hey")
	})*/

	// Start scraping on https://en.wikipedia.org
	c.Visit("https://www.exness-168.com/about_us/")
}