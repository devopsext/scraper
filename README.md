# Scraper

CLI utility to provide website cold scraping based on its tree.

It helps QA, Developers and DevOps to estimate time fractions for each request along side with response/status codes, content length and content type.

With Scraper json or yaml report you can easily understand consistency of website content. 

[![GoDoc](https://godoc.org/github.com/devopsext/scraper?status.svg)](https://godoc.org/github.com/devopsext/scraper)
[![build status](https://img.shields.io/travis/devopsext/scraper/master.svg?style=flat-square)](https://travis-ci.org/devopsext/scraper)

## Features

- Scrape whole website based on its tree
- Scan links, images, scripts, styles and emails
- Chrome rendering if browser option  is enabled
- Limit scraping by list of domains
- Track DNS, Connect, TLSHandshake, FirstByte, Download time
- Skip verify on wrong SSL certificates
- Provide response code, content length and type
- Follow redirects where it is needed
- Output to file: json or yaml

## Build

```sh
git clone https://github.com/devopsext/scraper.git
cd scraper/
go build
```

## Example

```sh
./scraper website --url https://devopsext.com --max-depth 3 --redirects \
                  --domains devopsext.com,www.devopsext.com \
                  --links --scripts --styles --images \
                  --output json --file scraper.json \
                  --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

## Get response codes

```sh
cat scraper.json | jq .StatusCodes
```

```json
{
  "200": 17
}
```

## Usage

```
Website command

Usage:
  scraper website [flags]

Flags:
      --browser string      Website browser: chrome
      --domains strings     Website domains
      --emails              Website scrape emails
      --file string         Website output file
  -h, --help                help for website
      --images              Website scrape images
      --insecure            Website insecure skip verify
      --links               Website scrape links
      --max-body-size int   Website max body size (default 10485760)
      --max-depth int       Website max depth (default 1)
      --output string       Website output: json, yaml
      --redirects           Website follow redirects
      --scripts             Website scrape scripts
      --styles              Website scrape styles
      --url string          Website url
      --user-agent string   Website user agent

Global Flags:
      --log-format string     Log format: json, text, stdout (default "text")
      --log-level string      Log level: info, warn, error, debug, panic (default "info")
      --log-template string   Log template (default "{{.func}} [{{.line}}]: {{.msg}}")
```

## Environment variables

For containerization purpose all command switches have environment variables analogs.

- SCRAPER_WEBSITE_URL
- SCRAPER_WEBSITE_REDIRECTS
- SCRAPER_WEBSITE_LINKS
- SCRAPER_WEBSITE_IMAGES
- SCRAPER_WEBSITE_SCRIPTS
- SCRAPER_WEBSITE_STYLES
- SCRAPER_WEBSITE_EMAILS
- SCRAPER_WEBSITE_DOMAINS
- SCRAPER_WEBSITE_OUTPUT
- SCRAPER_WEBSITE_MAX_DEPTH
- SCRAPER_WEBSITE_USER_AGENT
- SCRAPER_WEBSITE_INSECURE
- SCRAPER_WEBSITE_MAX_BODY_SIZE
- SCRAPER_WEBSITE_BROWSER
- SCRAPER_WEBSITE_FILE
- SCRAPER_LOG_FORMAT
- SCRAPER_LOG_LEVEL
- SCRAPER_LOG_TEMPLATE
