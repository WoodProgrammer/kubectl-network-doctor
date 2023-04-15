package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type dnsMapList struct {
	results []dnsMap
}
type dnsMap struct {
	name               string
	addrList           []string
	resolutionDuration time.Duration
	err                bool
}

func readHostList(file string) []string {

	var valueList []string
	readFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	for _, line := range fileLines {
		valueList = append(valueList, line)
	}

	return valueList

}

func resolveAddr(addr string, resultMap *dnsMap) *dnsMap {
	t1 := time.Now()
	_, err := net.LookupIP(addr)

	timeElapsed := time.Since(t1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		resultMap.err = true
	} else {
		resultMap.err = false
	}

	resultMap.name = addr

	/*for _, ip := range ips {
		resultMap.addrList = append(resultMap.addrList, ip.String())
	}*/
	resultMap.resolutionDuration = time.Duration(timeElapsed.Milliseconds())

	return resultMap
}

func main() {
	dummyMap := dnsMap{}
	dummyMapList := dnsMapList{}

	addrList := readHostList("hosts.txt")

	for _, addr := range addrList {
		result := resolveAddr(addr, &dummyMap)
		dummyMapList.results = append(dummyMapList.results, *result)
	}
	fmt.Println(dummyMapList)
}
