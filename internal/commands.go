package internal

import (
	"errors"
	"log"
	"net"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingResult struct {
	Successful bool          `json:"successful"`
	Time       time.Duration `json:"time"`
	Packets    int           `json:"packets"`
}

type SystemInfo struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
}

// Ping the host and return result
func Ping(host string) (PingResult, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return PingResult{}, err
	}

	pinger.SetPrivileged(true)

	// Set the number of packets
	pinger.Count = 4

	// Set the timeout
	pinger.Timeout = 10 * time.Second

	// Run the ping
	start := time.Now()
	err = pinger.Run() // Blocks until finished
	elapsed := time.Since(start)

	stats := pinger.Statistics()

	if err != nil {
		log.Printf("Ping failed: %v", err)
		return PingResult{Successful: false, Time: elapsed, Packets: stats.PacketsRecv}, err
	}

	return PingResult{
		Successful: stats.PacketsRecv > 0,
		Time:       elapsed,
		Packets:    stats.PacketsRecv,
	}, nil
}

// Get system hostname and IP
func GetSystemInfo() (SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return SystemInfo{}, err
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return SystemInfo{}, err
	}

	var ipAddress string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipAddress = ipNet.IP.String()
				break
			}
		}
	}

	if ipAddress == "" {
		return SystemInfo{}, errors.New("no IP address found")
	}

	return SystemInfo{
		Hostname:  hostname,
		IPAddress: ipAddress,
	}, nil
}
