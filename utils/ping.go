package utils

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Ping(address string) (float64, error) {
	ipAddr, err := net.ResolveIPAddr("ip4", address)
	if err != nil {
		return 0, fmt.Errorf("Failed to resolve IP: %v", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		return 0, fmt.Errorf("Failed to create ICMP connection: %v", err)
	}
	defer conn.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("HELLO!"),
		},
	}

	messageBytes, err := message.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to marshal ICMP message: %v", err)
	}

	start := time.Now()
	if _, err := conn.WriteTo(messageBytes, ipAddr); err != nil {
		return 0, fmt.Errorf("Failed to send ICMP request: %v", err)
	}

	reply := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return 0, fmt.Errorf("Failed to read ICMP response: %v", err)
	}

	parsedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		return 0, fmt.Errorf("Failed to parse ICMP response: %v", err)
	}

	if parsedMessage.Type != ipv4.ICMPTypeEchoReply {
		return 0, fmt.Errorf("Unexpected message type: %v", parsedMessage.Type)
	}

	duration := time.Since(start)
	return duration.Seconds() * 1000, nil
}

func NetPing(ip string, protocol string) (float64, error) {
	ports := []string{"443", "80"}

	for _, port := range ports {
		address := net.JoinHostPort(ip, port)

		start := time.Now()
		dialer := net.Dialer{Timeout: 1 * time.Second}
		conn, err := dialer.Dial(protocol, address)
		if err == nil {
			defer conn.Close()
			duration := time.Since(start)
			return duration.Seconds() * 1000, nil
		}
	}

	return 0, fmt.Errorf("unable to connect to %s on default ports", ip)
}
