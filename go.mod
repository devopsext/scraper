module github.com/devopsext/scraper

go 1.13

require (
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/antchfx/htmlquery v1.2.1 // indirect
	github.com/antchfx/xmlquery v1.2.2 // indirect
	github.com/chromedp/chromedp v0.5.2

	github.com/devopsext/colly v1.2.1-0.20191224114515-9ab1fb7c0704
	github.com/devopsext/utils v0.0.3
	github.com/go-delve/delve v1.3.2 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gocolly/colly v1.2.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/prometheus/common v0.7.0
	github.com/sirupsen/logrus v1.4.2

	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 // indirect
	golang.org/x/tools/gopls v0.2.2 // indirect
)

replace github.com/devopsext/scraper => ../../src/scraper

replace github.com/devopsext/colly => ../../src/colly

replace github.com/devopsext/utils => ../../src/utils
