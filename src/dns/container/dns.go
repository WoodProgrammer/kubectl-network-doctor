package dn

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func checkFileConfiguration() {
	var fileList [2]string

	fileList[0] = "/etc/hosts"
	fileList[1] = "/etc/resolv.conf"

	for _, file := range fileList {
		dat, _ := os.ReadFile(file)
		fmt.Println("######### The file from the container is %s #############", file)
		fmt.Println(string(dat))
		fmt.Println("####################################################")
	}
}

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
