package website

import (
	"net/url"
	"time"
)

type WebsitePageTime struct {
	DNS          time.Duration
	Connect      time.Duration
	TLSHandshake time.Duration
	FirstByte    time.Duration
	Download     time.Duration
}

type WebsitePage struct {
	URL        *url.URL
	Time       WebsitePageTime
	StatusCode int
	Parent     *WebsitePage
	Links      []string
	Children   []*WebsitePage
}
