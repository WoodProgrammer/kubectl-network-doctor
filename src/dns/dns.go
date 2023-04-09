package dns

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func dnsLookupMeasurement(addr string) {
	start := time.Now()
	ips, err := net.LookupIP(addr)

	log.Println("dial time: ", time.Since(start))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	}
	for _, ip := range ips {
		fmt.Printf("%s. IN A %s\n", addr, ip.String())
	}
}
