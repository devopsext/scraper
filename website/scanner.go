package website

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/devopsext/colly"
	"github.com/devopsext/scraper/common"
	"github.com/devopsext/utils"
)

var scannerLog = utils.GetLog()

type WebsiteScanner struct {
	page    *WebsitePage
	options common.WebsiteOptions
}

func (ws *WebsiteScanner) findPage(parent *WebsitePage, url string, depth int, level int) *WebsitePage {

	if parent == nil {
		return nil
	}

	for _, item := range parent.Children {

		if level == depth {
			return item
		} else {
			return ws.findPage(item, url, depth, level+1)
		}
	}

	return nil
}

func (ws *WebsiteScanner) findParentByURL(url string, depth int) *WebsitePage {

	if ws.page != nil {

		return ws.findPage(ws.page, url, depth, 0)
	}
	return nil
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

		parent := ws.findParentByURL(r.URL.String(), r.Depth)
		current = &WebsitePage{
			URL:    r.URL,
			Parent: parent,
		}

		if ws.page == nil && parent == nil {
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
		}

		scannerLog.Debug("response => %d", response.StatusCode)
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

			l := len(via)
			if l > 0 {
				scannerLog.Debug("response => %d to %s", req.Response.StatusCode, via[l-1].URL.String())

				next := &WebsitePage{
					URL:    via[l-1].URL,
					Parent: current,
				}
				current.Children = append(current.Children, next)
				current = next
			}

			return nil
		}
	}

	skipErrors := []interface{}{colly.ErrAlreadyVisited, colly.ErrMaxDepth}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		scannerLog.Debug("link: %s", link)

		if current != nil {
			current.Links = append(current.Links, link)
		}

		if err := e.Request.Visit(link); err != nil {

			if err == colly.ErrAlreadyVisited {
				scannerLog.Debug("already visited")
			} else if !utils.Contains(skipErrors, err) {
				scannerLog.Error(err)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		scannerLog.Error(err)
	})

	c.OnScraped(func(r *colly.Response) {
		if ws.page != nil {

		}
	})

	if err := c.Visit(ws.options.URL); err != nil {
		scannerLog.Error(err)
	}
}

func NewWebsiteScanner(options common.WebsiteOptions) *WebsiteScanner {

	return &WebsiteScanner{
		options: options,
	}
}
