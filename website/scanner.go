package website

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	"github.com/devopsext/scraper/common"
	"github.com/devopsext/utils"
	"github.com/gocolly/colly"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v2"
)

var scannerLog = utils.GetLog()

type pageDepth struct {
	depth int
	page  *WebsitePage
}

type WebsiteScanner struct {
	pageDepths map[uint32]pageDepth
	page       *WebsitePage
	options    common.WebsiteOptions
}

type WebsiteScannerResult struct {
	Root        *WebsitePage
	StatusCodes map[int]int
}

func resultToJson(r WebsiteScannerResult) []byte {

	b, err := json.Marshal(r)
	if err != nil {
		scannerLog.Error(err)
	}
	return b
}

func resultToYaml(r WebsiteScannerResult) []byte {

	b, err := yaml.Marshal(r)
	if err != nil {
		scannerLog.Error(err)
	}
	return b
}

func chromeSession(response *colly.Response) {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	if err := chromedp.Run(ctx, chromedp.Navigate(response.Request.URL.String()), chromedp.InnerHTML("html", &res)); err != nil {
		scannerLog.Error(err)
		return
	}

	response.Body = []byte(res)
}

func (ws *WebsiteScanner) getParentByRequest(r *colly.Request) *WebsitePage {

	ID := r.ID - 1
	for {

		pageDepth, ok := ws.pageDepths[ID]
		if ok {
			if pageDepth.depth < r.Depth {
				return pageDepth.page
			}
		} else {
			if ID < 1 {
				return nil
			}
		}
		ID--
	}
}

func (ws *WebsiteScanner) getDomain(s string) string {

	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}

func (ws *WebsiteScanner) Start() {

	var start, connect, dns, tlsHandshake, firstByte time.Time

	c := colly.NewCollector(
		colly.MaxDepth(ws.options.MaxDepth),
	)

	t := http.DefaultTransport.(*http.Transport)

	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: ws.options.Insecure}
	t.DisableKeepAlives = true

	c.Async = false
	c.AllowedDomains = ws.options.Domains

	domain := ws.getDomain(ws.options.URL)
	if !utils.Contains(c.AllowedDomains, domain) && !utils.IsEmpty(domain) {
		c.AllowedDomains = append(c.AllowedDomains, domain)
	}

	c.AllowURLRevisit = false
	c.MaxBodySize = ws.options.MaxBodySize

	if !utils.IsEmpty(ws.options.UserAgent) {
		c.UserAgent = ws.options.UserAgent
	}

	var current *WebsitePage
	codes := make(map[int]int)

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			if current != nil {
				current.Time.DNS = time.Since(dns)
			}
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			if current != nil {
				current.Time.Connect = time.Since(connect)
			}
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			if current != nil {
				current.Time.TLSHandshake = time.Since(tlsHandshake)
			}
		},

		GotFirstResponseByte: func() {
			firstByte = time.Now()
			if current != nil {
				current.Time.FirstByte = time.Since(start)
			}
		},
	}

	c.GetClientTrace = func(req *http.Request) *httptrace.ClientTrace {
		return trace
	}

	c.OnRequest(func(r *colly.Request) {

		start = time.Now()

		parent := ws.getParentByRequest(r)

		current = &WebsitePage{
			URL:    r.URL.String(),
			Parent: parent,
		}

		if ws.page == nil {
			ws.page = current
		}

		scannerLog.Debug("request => %s", r.URL.String())
	})

	c.OnResponse(func(response *colly.Response) {

		codes[response.StatusCode] = codes[response.StatusCode] + 1

		if current != nil {
			current.Time.Download = time.Since(firstByte)
			current.StatusCode = response.StatusCode
			current.Length = len(response.Body)

			contentType := response.Headers.Get("Content-Type")
			if !utils.IsEmpty(contentType) {
				current.Type = contentType
			}

			if current.Parent != nil {
				current.Parent.Children = append(current.Parent.Children, current)
			}

			ws.pageDepths[response.Request.ID] = pageDepth{
				depth: response.Request.Depth,
				page:  current,
			}
		}

		scannerLog.Debug("response %d", response.StatusCode)
		//scannerLog.Debug("response %d: %s", response.StatusCode, string(response.Body[:]))

		browser := strings.ToLower(ws.options.Browser)
		switch browser {
		case "chrome":
			chromeSession(response)
		}
	})

	if ws.options.Redirects {

		c.RedirectHandler = func(req *http.Request, via []*http.Request) error {

			codes[req.Response.StatusCode] = codes[req.Response.StatusCode] + 1

			if current != nil {
				current.Time.Download = time.Since(firstByte)
				current.StatusCode = req.Response.StatusCode

				if current.Parent != nil {
					current.Parent.Children = append(current.Parent.Children, current)
				}
			}

			scannerLog.Debug("response => %d to %s", req.Response.StatusCode, req.URL.String())

			next := &WebsitePage{
				URL:    req.URL.String(),
				Parent: current,
			}
			current = next

			return nil
		}
	}

	skipErrors := []interface{}{colly.ErrAlreadyVisited, colly.ErrMaxDepth, colly.ErrMissingURL, colly.ErrForbiddenDomain}

	if ws.options.Links {

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {

			link := e.Attr("href")
			if !utils.IsEmpty(link) {

				scannerLog.Debug("link: %s", link)

				if current != nil && !utils.Contains(current.Links, link) {
					current.Links = append(current.Links, link)
				}

				old := current
				if err := e.Request.Visit(link); err != nil {

					if !utils.Contains(skipErrors, err) {
						scannerLog.Error(err)
					}
				}
				current = old
			}
		})
	}

	if ws.options.Images {

		c.OnHTML("img[src]", func(e *colly.HTMLElement) {

			image := e.Attr("src")
			if !utils.IsEmpty(image) {

				scannerLog.Debug("image: %s", image)

				if current != nil && !utils.Contains(current.Images, image) {
					current.Images = append(current.Images, image)
				}

				old := current
				if err := e.Request.Visit(image); err != nil {

					if !utils.Contains(skipErrors, err) {
						scannerLog.Error(err)
					}
				}
				current = old
			}
		})
	}

	if ws.options.Scripts {

		c.OnHTML("script[src]", func(e *colly.HTMLElement) {

			script := e.Attr("src")
			if !utils.IsEmpty(script) {

				scannerLog.Debug("script: %s", script)

				if current != nil && !utils.Contains(current.Scripts, script) {
					current.Scripts = append(current.Scripts, script)
				}

				old := current
				if err := e.Request.Visit(script); err != nil {

					if !utils.Contains(skipErrors, err) {
						scannerLog.Error(err)
					}
				}
				current = old
			}
		})
	}

	if ws.options.Styles {

		c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {

			style := e.Attr("href")
			if !utils.IsEmpty(style) {

				scannerLog.Debug("style: %s", style)

				if current != nil && !utils.Contains(current.Styles, style) {
					current.Styles = append(current.Styles, style)
				}

				old := current
				if err := e.Request.Visit(style); err != nil {

					if !utils.Contains(skipErrors, err) {
						scannerLog.Error(err)
					}
				}
				current = old
			}
		})
	}

	c.OnError(func(r *colly.Response, err error) {
		scannerLog.Error(err)
	})

	c.OnScraped(func(r *colly.Response) {

		if current != nil && ws.options.Emails {
			parseEmails(r.Body, &current.Emails)
		}
	})

	if err := c.Visit(ws.options.URL); err != nil {
		scannerLog.Error(err)
	} else {

		c.Wait()

		output := strings.ToLower(ws.options.Output)
		result := WebsiteScannerResult{
			Root:        ws.page,
			StatusCodes: codes,
		}
		var b []byte

		switch output {
		case "json":
			b = resultToJson(result)
		case "yaml":
			b = resultToYaml(result)
		}

		if utils.IsEmpty(ws.options.File) {
			fmt.Println(string(b))
		} else {
			err := ioutil.WriteFile(ws.options.File, b, 0644)
			if err != nil {
				scannerLog.Error(err)
			}
		}
	}
}

func NewWebsiteScanner(options common.WebsiteOptions) *WebsiteScanner {

	return &WebsiteScanner{
		pageDepths: make(map[uint32]pageDepth),
		options:    options,
	}
}
