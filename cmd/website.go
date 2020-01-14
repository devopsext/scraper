package cmd

import (
	"reflect"
	"strings"

	"github.com/devopsext/scraper/common"
	"github.com/devopsext/scraper/website"
	"github.com/devopsext/utils"
	"github.com/spf13/cobra"
)

var websiteEnv = utils.GetEnvironment()
var websiteLog = utils.GetLog()

var websiteOpts = common.WebsiteOptions{

	URL: websiteEnv.Get("SCRAPER_WEBSITE_URL", "").(string),
	//Silent:      websiteEnv.Get("SCRAPER_WEBSITE_SILENT", false).(bool),
	Redirects:   websiteEnv.Get("SCRAPER_WEBSITE_REDIRECTS", false).(bool),
	Links:       websiteEnv.Get("SCRAPER_WEBSITE_LINKS", false).(bool),
	Images:      websiteEnv.Get("SCRAPER_WEBSITE_IMAGES", false).(bool),
	Scripts:     websiteEnv.Get("SCRAPER_WEBSITE_SCRIPTS", false).(bool),
	Styles:      websiteEnv.Get("SCRAPER_WEBSITE_STYLES", false).(bool),
	Emails:      websiteEnv.Get("SCRAPER_WEBSITE_EMAILS", false).(bool),
	Domains:     strings.Split(websiteEnv.Get("SCRAPER_WEBSITE_DOMAINS", "").(string), ","),
	Output:      websiteEnv.Get("SCRAPER_WEBSITE_OUTPUT", "").(string),
	MaxDepth:    websiteEnv.Get("SCRAPER_WEBSITE_MAX_DEPTH", 1).(int),
	UserAgent:   websiteEnv.Get("SCRAPER_WEBSITE_USER_AGENT", "").(string),
	Insecure:    websiteEnv.Get("SCRAPER_WEBSITE_INSECURE", false).(bool),
	MaxBodySize: websiteEnv.Get("SCRAPER_WEBSITE_MAX_BODY_SIZE", 10*1024*1024).(int),
	Browser:     websiteEnv.Get("SCRAPER_WEBSITE_BROWSER", "").(string),
	File:        websiteEnv.Get("SCRAPER_WEBSITE_FILE", "").(string),
}

func GetWebsiteCmd() *cobra.Command {

	rootCmd := cobra.Command{
		Use:   "website",
		Short: "Website command",
		Run: func(cmd *cobra.Command, args []string) {

			scanner := website.NewWebsiteScanner(websiteOpts)
			if reflect.ValueOf(scanner).IsNil() {
				websiteLog.Panic("Website scanner is invalid. Terminating...")
			}

			scanner.Start()
		},
	}

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&websiteOpts.URL, "url", websiteOpts.URL, "Website url")
	//flags.BoolVar(&websiteOpts.Silent, "silent", websiteOpts.Silent, "Website silency")
	flags.BoolVar(&websiteOpts.Redirects, "redirects", websiteOpts.Redirects, "Website follow redirects")
	flags.BoolVar(&websiteOpts.Links, "links", websiteOpts.Links, "Website scrape links")
	flags.BoolVar(&websiteOpts.Images, "images", websiteOpts.Images, "Website scrape images")
	flags.BoolVar(&websiteOpts.Scripts, "scripts", websiteOpts.Scripts, "Website scrape scripts")
	flags.BoolVar(&websiteOpts.Styles, "styles", websiteOpts.Styles, "Website scrape styles")
	flags.BoolVar(&websiteOpts.Emails, "emails", websiteOpts.Emails, "Website scrape emails")
	flags.StringSliceVar(&websiteOpts.Domains, "domains", websiteOpts.Domains, "Website domains")
	flags.StringVar(&websiteOpts.Output, "output", websiteOpts.Output, "Website output: json, yaml")
	flags.IntVar(&websiteOpts.MaxDepth, "max-depth", websiteOpts.MaxDepth, "Website max depth")
	flags.StringVar(&websiteOpts.UserAgent, "user-agent", websiteOpts.UserAgent, "Website user agent")
	flags.BoolVar(&websiteOpts.Insecure, "insecure", websiteOpts.Insecure, "Website insecure skip verify")
	flags.IntVar(&websiteOpts.MaxBodySize, "max-body-size", websiteOpts.MaxBodySize, "Website max body size")
	flags.StringVar(&websiteOpts.Browser, "browser", websiteOpts.Browser, "Website browser: chrome")
	flags.StringVar(&websiteOpts.File, "file", websiteOpts.File, "Website output file")

	return &rootCmd
}
