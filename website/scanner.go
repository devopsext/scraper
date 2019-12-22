package website

import (


	"github.com/devopsext/scraper/common"
)

type WebsiteScanner struct {
	options common.WebsiteOptions
}

func (ws *WebsiteScanner) getUrl(s string) string {

	return ""
}

func NewWebsiteScanner(options common.WebsiteOptions) *WebsiteScanner {

	return &WebsiteScanner{
		options: options,
	}
}
