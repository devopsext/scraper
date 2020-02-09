package common

type WebsiteScanner interface {
	Start()
}

type WebsiteOptions struct {
	URL         string
	Redirects   bool
	Links       bool
	Images      bool
	Scripts     bool
	Styles      bool
	Emails      bool
	Domains     []string
	Output      string
	MaxDepth    int
	UserAgent   string
	Insecure    bool
	MaxBodySize int
	Browser     string
	File        string
}
