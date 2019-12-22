package common

type WebsiteScanner interface {
}

type WebsiteOptions struct {
	Url       string
	Silent    bool
	Redirects bool
	Domains   []string
}
