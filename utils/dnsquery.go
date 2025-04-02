package utils

import (
	"net"
)

func DnsQuery(domain string) []string {
	var ips []string

	query, _ := net.LookupIP(domain)
	for _, ip := range query {
		if ipv4 := ip.To4(); ipv4 != nil {
			ips = append(ips, ipv4.String())
		}
	}

	return ips
}
