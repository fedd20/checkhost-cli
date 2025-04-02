package utils

import (
	"checkhost-cli/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func LoadConfig() (*viper.Viper, error) {
	log := logger.GetLoggerInstance(false)
	v := viper.New()

	v.SetConfigName(".checkhost-cli")
	v.SetConfigType("json")
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")

	v.SetDefault("CLOUDFLARE_API_URL", "https://get-location-from-ip.imperialwool.workers.dev/")
	v.SetDefault("SHOULD_SHOW_COUNTRY_FLAG", false)
	err, data := fetchCountryData()
	if err == nil {
		v.SetDefault("COUNTRY_DATA", data)
	} else {
		log.Warn().Err(err).Msg("Can't retrieve country data. It's better to recreate config later.")
		v.SetDefault("COUNTRY_DATA", map[string]any{})
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := v.SafeWriteConfig(); err != nil {
				return nil, fmt.Errorf("Error creating config file: %v", err)
			}
		} else {
			return nil, fmt.Errorf("Error reading config file: %v", err)
		}
	}

	return v, nil
}

func fetchCountryData() (error, map[string]any) {
	// JSON that contains data in format {"code": {"name": "Country Name", "emoji": "Country Emoji"}}
	url := "https://gist.githubusercontent.com/imperialwool/ecc7d5c1cdbe1cad1f1d961000c61897/raw/f797c6a20bf0176994e751d39feca50b065ac7e0/country-code-to-name-and-emoji.json"

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	var jsonData map[string]any
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return err, nil
	}

	return nil, jsonData
}
