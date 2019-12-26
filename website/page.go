package website

import (
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
	URL        string
	Time       WebsitePageTime
	StatusCode int
	Parent     *WebsitePage `json:"-" yaml:"-"`
	Links      []string
	Children   []*WebsitePage
}
