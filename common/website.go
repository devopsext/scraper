package common

type WebsiteScanner interface {
	Start()
}

type WebsiteOptions struct {
	URL         string
	Silent      bool
	Redirects   bool
	Domains     []string
	Output      string
	MaxDepth    int
	UserAgent   string
	Insecure    bool
	MaxBodySize int
	Browser     string
	File        string
}
