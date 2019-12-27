module github.com/devopsext/scraper

go 1.13

require (
	github.com/chromedp/chromedp v0.5.2

	github.com/devopsext/utils v0.0.3
	github.com/go-delve/delve v1.3.2 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/prometheus/common v0.7.0
	github.com/sirupsen/logrus v1.4.2

	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 // indirect
	golang.org/x/tools/gopls v0.2.2 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/gocolly/colly => github.com/devopsext/colly v1.2.1-0.20191227100724-341fb938e4bb

//replace github.com/devopsext/utils => ../../src/utils

//replace github.com/devopsext/scraper => ../../src/scraper
//replace github.com/gocolly/colly => ../../src/colly
