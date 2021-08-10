package main

import (
	"bufio"
	"flag"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var defaultTimeout = time.Second

func main() {
	inputFile := flag.String("i", "input.txt", "IP or domain names to scan")
	outputFile := flag.String("o", "output.txt", "Result list")
	timeout := flag.Int("t", 1, "Timeout of port scan")
	flag.Parse()

	if *timeout > 1 {
		defaultTimeout = time.Duration(*timeout) * time.Second
	}

	if *inputFile == "" {
		println(os.Args[0])
		println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	i, err := os.Open(*inputFile)
	if err != nil {
		println("input file:", err)
		os.Exit(1)
	}
	r := bufio.NewReader(i)

	o, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		println("output file:", err)
	}

	hosts := make([]string, 0)
	for {
		l, _, err := r.ReadLine()
		if err != nil || err == io.EOF {
			break
		}
		hosts = append(hosts, string(l))
	}

	for _, h := range hosts {
		println("host:", h)
		results := probeHost(h)
		println("found opened ports:", len(results))
		for _, p := range results {
			if _, err := o.WriteString(p + "\n"); err != nil {
				println("save:", err)
			}
		}
	}

	println("done")
}

func probeHost(host string) []string {
	probeResults := make([]string, 0)
	wg := sync.WaitGroup{}
	for i := 1; i <= 100_000; i++ {
		go func(port int) {
			wg.Add(1)
			defer wg.Done()

			address := net.JoinHostPort(host, strconv.Itoa(port))
			if isAddressActive(address) {
				probeResults = append(probeResults, address)
			}
		}(i)
	}
	wg.Wait()
	return probeResults
}

func isAddressActive(address string) bool {
	conn, err := net.DialTimeout("tcp", address, defaultTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
