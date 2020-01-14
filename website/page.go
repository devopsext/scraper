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
	Length     int
	Type       string
	Parent     *WebsitePage   `json:"-" yaml:"-"`
	Links      []string       `json:",omitempty" yaml:",omitempty"`
	Images     []string       `json:",omitempty" yaml:",omitempty"`
	Scripts    []string       `json:",omitempty" yaml:",omitempty"`
	Styles     []string       `json:",omitempty" yaml:",omitempty"`
	Emails     []string       `json:",omitempty" yaml:",omitempty"`
	Children   []*WebsitePage `json:",omitempty" yaml:",omitempty"`
}
