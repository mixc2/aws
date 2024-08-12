// By:@WilerHttpCC  转发请不要修改信息
// 免费全球高质量代理 频道@WilerP2C_proxylist
// Proxy_pool  代理池https://proxy.jsjiuah.life/
// 付费代理池仅需30U 全球代理存活8000+

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var awsDomainPatterns = []string{
	".compute-1.amazonaws.com",
	".compute.amazonaws.com",
	".compute-1.amazonaws.com.cn",
	".compute.amazonaws.com.cn",
	".amazonaws.com",
}

func awsIP(ip string) bool {
	addrs, err := net.LookupAddr(ip)
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		for _, pattern := range awsDomainPatterns {
			if strings.Contains(addr, pattern) {
				return true
			}
		}
	}
	return false
}

func filterProxy(proxy string) (string, bool) {
	ipPort := strings.Split(proxy, ":")
	ip := ipPort[0]

	if !awsIP(ip) {
		return proxy, true
	} else {
		fmt.Printf("Proxy %s AWS filtered\n", proxy)
		return "", false
	}
}

func filtrationProxy(inputFile, outputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Unable to open input file:", err)
		return
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Unable to create output file:", err)
		return
	}
	defer output.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			if filteredProxy, ok := filterProxy(proxy); ok {
				mu.Lock()
				output.WriteString(filteredProxy + "\n")
				mu.Unlock()
			}
		}(scanner.Text())
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading file:", err)
	}

	fmt.Printf("Saved as %s\n", outputFile)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run 1.go input_file output_file")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	filtrationProxy(inputFile, outputFile)
}
