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

	URL:       websiteEnv.Get("SCRAPER_WEBSITE_URL", "").(string),
	Silent:    websiteEnv.Get("SCRAPER_WEBSITE_SILENT", false).(bool),
	Redirects: websiteEnv.Get("SCRAPER_WEBSITE_REDIRECTS", false).(bool),
	Domains:   strings.Split(websiteEnv.Get("SCRAPER_WEBSITE_DOMAINS", "ya.ru").(string), ","),
	Output:    websiteEnv.Get("SCRAPER_WEBSITE_OUTPUT", "json").(string),
	MaxDepth:  websiteEnv.Get("SCRAPER_WEBSITE_MAX_DEPTH", 1).(int),
	UserAgent: websiteEnv.Get("SCRAPER_WEBSITE_USER_AGENT", "").(string),
	Insecure:  websiteEnv.Get("SCRAPER_WEBSITE_INSECURE", false).(bool),
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
	flags.BoolVar(&websiteOpts.Silent, "silent", websiteOpts.Silent, "Website silency")
	flags.BoolVar(&websiteOpts.Redirects, "redirects", websiteOpts.Redirects, "Website follow redirects")
	flags.StringSliceVar(&websiteOpts.Domains, "domains", websiteOpts.Domains, "Website domains")
	flags.StringVar(&websiteOpts.Output, "output", websiteOpts.Output, "Website output")
	flags.IntVar(&websiteOpts.MaxDepth, "max-depth", websiteOpts.MaxDepth, "Website max depth")
	flags.StringVar(&websiteOpts.UserAgent, "user-agent", websiteOpts.UserAgent, "Website user agent")
	flags.BoolVar(&websiteOpts.Insecure, "insecure", websiteOpts.Insecure, "Website insecure skip verify")

	return &rootCmd
}
