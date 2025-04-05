package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started on port %d", addr.Port)

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	rand.Seed(time.Now().UnixNano())

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("error ", err)
			continue
		}

		msg := string(buf[:n])
		log.Printf("Received: %s", msg)

		if rand.Intn(100) < 20 {
			log.Printf("Packet lost: %s", msg)
			continue
		}

		response := strings.ToUpper(msg)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			log.Println("error:", err)
		} else {
			log.Printf("Sent: %s", response)
		}
	}
}
