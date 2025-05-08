package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	firstTask()
	secondTask()
}

func firstTask() {
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil || ipNet.IP.IsLoopback() {
				continue
			}

			fmt.Printf("IP-адрес -  %s\n", ipNet.IP)
			fmt.Printf("Маска    -   %s\n\n", net.IP(ipNet.Mask))
		}
	}
}

func secondTask() {
	if len(os.Args) != 4 {
		fmt.Println("ip start end")
		os.Exit(1)
	}

	ip := os.Args[1]
	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	end, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if start > end {
		start, end = end, start
	}

	var wg sync.WaitGroup
	ports := make(chan int)

	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if isPortAvailable(ip, p) {
				ports <- p
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ports)
	}()

	fmt.Println("доступны - ")
	for port := range ports {
		fmt.Print(port, " ")
	}
}

func isPortAvailable(ip string, port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, fmt.Sprintf("%d", port)), time.Second)
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}
