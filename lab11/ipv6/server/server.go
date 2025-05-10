package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp6", "[::1]:12345")
	if err != nil {
		log.Println("ошибка запуска - ", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("ошибка подлкючения клиента - ", err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for {
				message, err := reader.ReadString('\n')
				if err != nil {
					log.Println("клиент ушел - ", err)
					return
				}
				resp := strings.ToUpper(message)
				fmt.Printf("былo - %s", message)
				fmt.Printf("стало - %s", resp)
				_, err = conn.Write([]byte(resp))
				if err != nil {
					log.Println("ошибка отправки", err)
					return
				}
			}
		}(conn)
	}
}
