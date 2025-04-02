package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IpApiResponse struct {
	Ip        string `json:"query"`
	Asn       string `json:"as"`
	Country   string `json:"countryCode"`
	Timezone  string `json:"timezone"`
	City      string `json:"city"`
	Region    string `json:"regionName"`
	IsMobile  bool   `json:"mobile,omitempty"`
	IsProxy   bool   `json:"proxy,omitempty"`
	IsHosting bool   `json:"hosting,omitempty"`
}

func IpApiRequest(ip string) (*IpApiResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,countryCode,regionName,city,timezone,as,mobile,proxy,hosting,query", ip))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var data IpApiResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
