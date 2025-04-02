package apis

import (
	"encoding/json"
	"net/http"
)

const (
	CloudflareApiUrl = "https://get-location-from-ip.imperialwool.workers.dev/"
	// [WARNING]
	// Please note that this API is limited since its hosted for free.
	// It can be called for 100k times per day only. Every request by everyone is counted.
	// You should fork it, host it yourself at Cloudflare Workers and change the URL to your own.
)

type CloudflareApiResponse struct {
	Ip        string `json:"ip"`
	Asn       string `json:"asn"`
	Country   string `json:"country"`
	Timezone  string `json:"timezone"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Postcode  string `json:"postcode"`
}

func CloudflareApiRequest() (*CloudflareApiResponse, error) {
	resp, err := http.Get(CloudflareApiUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var data CloudflareApiResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
