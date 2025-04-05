package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", os.Args[1]+":"+os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 1; i <= 10; i++ {
		msg := fmt.Sprintf("Ping %d %s", i, time.Now().Format("15:04:05"))
		start := time.Now()

		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Println("error:", err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Request %d timed out\n", i)
			} else {
				log.Println("error:", err)
			}
			continue
		}
		fmt.Printf("%s, RTT: %v\n", string(buf[:n]), time.Since(start))
	}
}
