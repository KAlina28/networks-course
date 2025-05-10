package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp6", "[::1]:12345")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("ошибка ввода - ", err)
			break
		}
		_, err = conn.Write([]byte(text))
		if err != nil {
			log.Println("ошибка отправки - ", err)
			break
		}
		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("ошибка чтения - ", err)
			break
		}
		fmt.Printf("ответ сервера - %s", resp)
	}
}
