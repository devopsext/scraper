package website

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/devopsext/colly"
	"github.com/devopsext/scraper/common"
	"github.com/devopsext/utils"

	"github.com/go-yaml/yaml"
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
	if len(c.AllowedDomains) == 0 {
		c.AllowedDomains = []string{ws.options.URL}
	}
	c.AllowURLRevisit = false
	c.MaxBodySize = ws.options.MaxBodySize

	if !utils.IsEmpty(ws.options.UserAgent) {
		c.UserAgent = ws.options.UserAgent
	}

	var current *WebsitePage

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

		if current != nil {
			current.Time.Download = time.Since(firstByte)
			current.StatusCode = response.StatusCode

			if current.Parent != nil {
				current.Parent.Children = append(current.Parent.Children, current)
			}

			ws.pageDepths[response.Request.ID] = pageDepth{
				depth: response.Request.Depth,
				page:  current,
			}

		}

		//scannerLog.Debug("response %d: %s", response.StatusCode, string(response.Body[:]))
		scannerLog.Debug("response %d", response.StatusCode)
	})

	if ws.options.Redirects {

		c.RedirectHandler = func(req *http.Request, via []*http.Request) error {

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

	skipErrors := []interface{}{colly.ErrAlreadyVisited, colly.ErrMaxDepth, colly.ErrMissingURL}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		scannerLog.Debug("link: %s", link)

		needToVisit := true

		if current != nil {
			needToVisit = !utils.Contains(current.Links, link)
			if needToVisit {
				current.Links = append(current.Links, link)
			}
		}

		if needToVisit {
			if err := e.Request.Visit(link); err != nil {

				if err == colly.ErrAlreadyVisited {
					scannerLog.Debug("already visited")
				} else if !utils.Contains(skipErrors, err) {
					//scannerLog.Error(err)
				}
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		scannerLog.Error(err)
	})

	c.OnScraped(func(r *colly.Response) {

	})

	if err := c.Visit(ws.options.URL); err != nil {
		scannerLog.Error(err)
	} else {
		if ws.page != nil {

			output := strings.ToLower(ws.options.Output)
			if output == "json" {

				b, err := json.Marshal(ws.page)
				if err != nil {
					scannerLog.Error(err)
				}
				fmt.Println(string(b))

			} else if output == "yaml" {

				b, err := yaml.Marshal(ws.page)
				if err != nil {
					scannerLog.Error(err)
				}
				fmt.Println(string(b))
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
