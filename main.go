package main

import (
	"context"
	"flag"
	"fmt"
	"net"
    "math"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Println("DNS Query and TCPing Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  dnstool -d <domain> -s <dns_servers> [-t <query_type>]")
	fmt.Println("\nParameters:")
	fmt.Println("  -d  Domain to query (required)")
	fmt.Println("  -s  Comma-separated list of DNS servers (required)")
	fmt.Println("  -t  Query type: 4 for A (IPv4), 6 for AAAA (IPv6) (default: 4)")
	fmt.Println("\nExample:")
	fmt.Println("  dnstool -d www.example.com -s 1.1.1.1,8.8.8.8,223.5.5.5 -t 4")
	fmt.Println("\nDescription:")
	fmt.Println("  This tool performs DNS queries for a specified domain using the provided DNS servers.")
	fmt.Println("  It then conducts a TCPing test to the resolved IP address on port 80.")
	fmt.Println("  Results include DNS query time and TCPing latency for each server.")
}

func main() {
	domain := flag.String("d", "", "Domain to query")
	servers := flag.String("s", "", "Comma-separated list of DNS servers")
	queryType := flag.String("t", "4", "Query type: 4 for A, 6 for AAAA")

	flag.Usage = usage
	flag.Parse()

	if *domain == "" || *servers == "" {
		fmt.Println("Error: Missing required parameters")
		usage()
		os.Exit(1)
	}

	dnsServers := strings.Split(*servers, ",")

	var recordType uint16
	if *queryType == "6" {
		recordType = net.IPv6len
	} else {
		recordType = net.IPv4len
	}

	for _, server := range dnsServers {
		start := time.Now()
		ip, err := lookupIP(*domain, server, recordType)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("Error querying %s using %s: %v\n", *domain, server, err)
			continue
		}

		fmt.Printf("DNS Server: %s, Query time: %v\n", server, duration)
		fmt.Printf("IP: %s\n", ip)

		pingDuration, err := tcping(ip, 80)
		if err != nil {
			fmt.Printf("TCPing error: %v\n", err)
		} else {
			fmt.Printf("Average TCPing time: %.2fms\n", float64(pingDuration)/float64(time.Millisecond))
		}
		fmt.Println()
	}
}

func lookupIP(domain, server string, recordType uint16) (string, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", server+":53")
		},
	}

	ips, err := r.LookupIP(context.Background(), "ip", domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if len(ip) == int(recordType) {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no matching IP found")
}

func tcping(ip string, port int) (time.Duration, error) {
    var addr string
    if strings.Contains(ip, ":") {
        // IPv6 address
        addr = fmt.Sprintf("[%s]:%d", ip, port)
    } else {
        // IPv4 address
        addr = fmt.Sprintf("%s:%d", ip, port)
    }

    var totalDuration time.Duration
    attempts := 4

    for i := 0; i < attempts; i++ {
        start := time.Now()
        conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
        if err != nil {
            return 0, err
        }
        duration := time.Since(start)
        totalDuration += duration
        conn.Close()

        if i < attempts-1 {
            time.Sleep(600 * time.Millisecond) // 在每次测试之间稍作暂停
        }
    }

    averageDuration := totalDuration / time.Duration(attempts)
    roundedDuration := time.Duration(math.Round(float64(averageDuration)/float64(time.Millisecond)*100) / 100 * float64(time.Millisecond))

    return roundedDuration, nil
}