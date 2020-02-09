module github.com/devopsext/scraper

go 1.13

require (
	github.com/antchfx/htmlquery v1.2.2 // indirect
	github.com/antchfx/xmlquery v1.2.3 // indirect
	github.com/antchfx/xpath v1.1.4 // indirect
	github.com/chromedp/cdproto v0.0.0-20200119225551-e0b5a74c467f // indirect
	github.com/chromedp/chromedp v0.5.3

	github.com/devopsext/utils v0.0.3
	github.com/go-delve/delve v1.3.2 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lawzava/scrape v1.4.0
	github.com/prometheus/common v0.7.0
	github.com/sirupsen/logrus v1.4.2

	github.com/spf13/cobra v0.0.5
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200120151820-655fe14d7479 // indirect
	golang.org/x/tools/gopls v0.2.2 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/gocolly/colly => github.com/devopsext/colly v1.2.1-0.20191227100724-341fb938e4bb

//replace github.com/devopsext/utils => ../../src/utils

//replace github.com/devopsext/scraper => ../../src/scraper
//replace github.com/gocolly/colly => ../../src/colly
