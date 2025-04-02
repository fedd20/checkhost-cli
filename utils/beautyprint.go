package utils

import (
	"checkhost-cli/apis"
	"checkhost-cli/logger"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type CountryInfo struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

func getCountryInfo(config *viper.Viper, country string) CountryInfo {
	log := logger.GetLoggerInstance(false)
	countryInfo := CountryInfo{
		Name:  "Unknown",
		Emoji: "❓",
	}

	countryData := config.GetStringMap("country_data")
	if countryData == nil {
		return countryInfo
	}

	country = strings.ToLower(country)
	if data, ok := countryData[country]; ok {
		if countryMap, ok := data.(map[string]any); ok {
			if name, ok := countryMap["name"].(string); ok {
				countryInfo.Name = name
			} else {
				log.Warn().Str("country", country).Msg("Can't get `name` parameter for country. Not critical because we can fallback.")
			}
			if emoji, ok := countryMap["emoji"].(string); ok {
				countryInfo.Emoji = emoji
			} else {
				log.Warn().Str("country", country).Msg("Can't get `emoji` parameter for country. Not critical because we can fallback.")
			}
		}
	}

	return countryInfo
}

func compactPing(log *zerolog.Logger, ip string) float64 {
	ping, err := Ping(ip)
	if err != nil {
		log.Warn().Str("ip", ip).Err(err).Msg("Error pinging ICMP IP")
		ping, err = NetPing(ip, "tcp")
		if err != nil {
			log.Warn().Str("ip", ip).Err(err).Msg("Error pinging TCP IP")
			ping, err = NetPing(ip, "udp")
			if err != nil {
				log.Warn().Str("ip", ip).Err(err).Msg("Error pinging UDP IP")
				ping = 0
			}
		}
	}

	return ping
}

func BeautyPrint(config *viper.Viper, data any, isLocal bool) {
	log := logger.GetLoggerInstance(false)
	strings := []string{
		"\n█▀▀ █░█ █▀▀ █▀▀ █▄▀ █░█ █▀█ █▀ ▀█▀\n",
		"█▄▄ █▀█ ██▄ █▄▄ █░█ █▀█ █▄█ ▄█ ░█░\n\n",
	}
	showCountryFlag := config.GetBool("SHOULD_SHOW_COUNTRY_FLAG")

	switch typedData := data.(type) {
	case *apis.CloudflareApiResponse:
		if typedData == nil {
			log.Fatal().Msg("No data recieved. Probably request issue.")
		}
		countryInfo := getCountryInfo(config, typedData.Country)

		if showCountryFlag {
			strings = append(strings, fmt.Sprintf("%s %s\t💻 IP: %s\n", countryInfo.Emoji, countryInfo.Name, typedData.Ip))
		} else {
			strings = append(strings, fmt.Sprintf("🌌 %s\t💻 IP: %s\n", countryInfo.Name, typedData.Ip))
		}
		strings = append(strings, fmt.Sprintf("⌚ %s\t🌐 %s\n", typedData.Timezone, typedData.Asn))
		strings = append(strings, fmt.Sprintf("🧭 %s, %s\t📶 Ping: %f ms\n", typedData.Region, typedData.City, 0.00))
		strings = append(strings, fmt.Sprintf("📌 %s, %s\n", typedData.Latitude, typedData.Longitude))
		strings = append(strings, fmt.Sprintf("📭 %s\n\n", typedData.Postcode))
	case *apis.IpApiResponse:
		if typedData == nil {
			log.Fatal().Msg("No data recieved. Probably request issue.")
		}
		countryInfo := getCountryInfo(config, typedData.Country)
		var ping float64
		if isLocal {
			ping = 0
		} else {
			ping = compactPing(log, typedData.Ip)
		}

		if showCountryFlag {
			strings = append(strings, fmt.Sprintf("%s %s\t💻 IP: %s\n", countryInfo.Emoji, countryInfo.Name, typedData.Ip))
		} else {
			strings = append(strings, fmt.Sprintf("🌌 %s\t💻 IP: %s\n", countryInfo.Name, typedData.Ip))
		}
		strings = append(strings, fmt.Sprintf("⌚ %s\t🌐 %s\n", typedData.Timezone, typedData.Asn))
		strings = append(strings, fmt.Sprintf("🧭 %s, %s\t📶 Ping: %f ms\n\n", typedData.Region, typedData.City, ping))
		if typedData.IsMobile {
			strings = append(strings, fmt.Sprintf("📱 Mobile or residential IP.\n\n"))
		}
		if typedData.IsProxy {
			strings = append(strings, fmt.Sprintf("🎭 Proxy IP.\n\n"))
		}
		if typedData.IsHosting {
			strings = append(strings, fmt.Sprintf("🏢 Hosting IP.\n\n"))
		}
	case *apis.IpInfoResponse:
		if typedData == nil {
			log.Fatal().Msg("No data recieved. Probably request issue.")
		}
		countryInfo := getCountryInfo(config, typedData.Country)
		var ping float64
		if isLocal {
			ping = 0
		} else {
			ping = compactPing(log, typedData.Ip)
		}

		if showCountryFlag {
			strings = append(strings, fmt.Sprintf("%s %s\t💻 IP: %s\n", countryInfo.Emoji, countryInfo.Name, typedData.Ip))
		} else {
			strings = append(strings, fmt.Sprintf("🌌 %s\t💻 IP: %s\n", countryInfo.Name, typedData.Ip))
		}
		strings = append(strings, fmt.Sprintf("⌚ %s\t🌐 %s\n", typedData.Timezone, typedData.Asn))
		strings = append(strings, fmt.Sprintf("🧭 %s, %s\t📶 Ping: %f ms\n\n", typedData.Region, typedData.City, ping))
	default:
		log.Fatal().Msg("API service that user choose is not exist or not implemented")
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 4, ' ', 0)
	for _, l := range strings {
		tw.Write([]byte(l))
	}
	tw.Flush()
}
