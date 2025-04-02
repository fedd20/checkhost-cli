package cmd

import (
	"checkhost-cli/apis"
	"checkhost-cli/logger"
	"checkhost-cli/utils"
	"net"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var (
	source string
	debug  bool
)

var rootCmd = &cobra.Command{
	Use:   "checkhost",
	Short: "checkhost is a cli tool to get basic info about domain or IP",
	Long:  "checkhost is a cli tool to get basic info about domain or IP - location, ping, ASN, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.GetLoggerInstance(debug)

		var lookupIp string
		if len(args) == 0 {
			log.Info().Msg("Looking up public IP address of machine...")
			lookupIp = ""
		} else {
			lookupIp = args[0]
			if lookupIp == "" {
				log.Fatal().Msg("Input cannot be empty.")
			}

			if net.ParseIP(lookupIp) == nil {
				url, err := url.Parse(lookupIp)
				cobra.CheckErr(err)

				domain := url.String()
				if domain == "" {
					log.Fatal().Str("lookupIp", lookupIp).Err(err).Msg("Can't parse URL")
				}

				ips := utils.DnsQuery(domain)
				if len(ips) == 0 {
					log.Fatal().Msg("No IP address found for provided URL.")
				}
				lookupIp = ips[0]
			}
		}

		config, err := utils.LoadConfig()
		if err != nil {
			log.Fatal().Err(err)
		}

		switch source {
		case "ipinfo":
			result, err := apis.IpInfoRequest(lookupIp)
			if err != nil {
				log.Fatal().Err(err)
			}
			utils.BeautyPrint(config, result, lookupIp == "")
		case "ipapi":
			result, err := apis.IpApiRequest(lookupIp)
			if err != nil {
				log.Fatal().Err(err)
			}
			utils.BeautyPrint(config, result, lookupIp == "")
		case "cloudflare":
			if lookupIp != "" {
				log.Warn().Msg("Please note that Cloudflare API can only lookup your IP address.")
			}
			result, err := apis.CloudflareApiRequest()
			if err != nil {
				log.Fatal().Err(err)
			}
			utils.BeautyPrint(config, result, lookupIp == "")
		default:
			log.Fatal().Str("source", source).Msg("Invalid source")
		}

		os.Exit(0)
	},
}

func Execute() error {
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show debug logs")
	rootCmd.PersistentFlags().StringVar(&source, "source", "ipapi", "Choose a source to retrive IP info.")
	return rootCmd.Execute()
}
