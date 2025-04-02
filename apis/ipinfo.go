package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IpInfoResponse struct {
	Ip       string `json:"ip"`
	Asn      string `json:"org"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
	City     string `json:"city"`
	Region   string `json:"region"`
}

func IpInfoRequest(ip string) (*IpInfoResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s", ip))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var data IpInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
