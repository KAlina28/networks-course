package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	from := flag.String("from", "", "отправитель")
	to := flag.String("to", "", "получатель")
	body := flag.String("body", "", "какой-то текст")
	attachment := flag.String("attach", "", "путь к файлу с картинкой") // для задания 3
	host := flag.String("host", "smtp.gmail.com", "SMTP сервер")
	username := flag.String("user", "", "логин")
	password := flag.String("pass", "", "пароль")

	flag.Parse()

	if *from == "" || *to == "" || *username == "" || *password == "" {
		fmt.Println("-from, -to, -user, -pass")
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.DialTimeout("tcp", *host+":587", 15*time.Second)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	if _, err := readingResp(reader); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	sendCommand(writer, reader, "EHLO client.example.com")
	sendCommand(writer, reader, "STARTTLS")

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: *host,
	})
	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	reader = bufio.NewReader(tlsConn)
	writer = bufio.NewWriter(tlsConn)
	sendCommand(writer, reader, "EHLO client.example.com")
	sendCommand(writer, reader, "AUTH LOGIN")
	sendCommand(writer, reader, base64.StdEncoding.EncodeToString([]byte(*username)))
	sendCommand(writer, reader, base64.StdEncoding.EncodeToString([]byte(*password)))
	sendCommand(writer, reader, fmt.Sprintf("MAIL FROM:<%s>", *from))
	sendCommand(writer, reader, fmt.Sprintf("RCPT TO:<%s>", *to))
	sendCommand(writer, reader, "DATA")
	writer.WriteString("MIME-Version: 1.0\r\n")
	writer.WriteString(fmt.Sprintf("From: %s\r\n", *from))
	writer.WriteString(fmt.Sprintf("To: %s\r\n", *to))

	if *attachment != "" {
		fileData, err := ioutil.ReadFile(*attachment)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		boundary := "boundary_" + fmt.Sprintf("%d", time.Now().UnixNano())
		writer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary))

		writer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		writer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
		writer.WriteString(*body + "\r\n")

		writer.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		writer.WriteString(fmt.Sprintf("Content-Type: %s\r\n", "image/jpeg"))
		writer.WriteString("Content-Transfer-Encoding: base64\r\n")
		writer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filepath.Base(*attachment)))
		encoded := base64.StdEncoding.EncodeToString(fileData)
		for i := 0; i < len(encoded); i += 76 {
			end := min(i+76, len(encoded))
			writer.WriteString(encoded[i:end] + "\r\n")
		}

		writer.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	} else {
		writer.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
		writer.WriteString(*body + "\r\n")
	}

	sendCommand(writer, reader, ".\r\n")
	sendCommand(writer, reader, "QUIT")
	fmt.Println("отправлено")
}

func sendCommand(writer *bufio.Writer, reader *bufio.Reader, cmd string) {
	fmt.Printf("-> %s\n", cmd)
	if _, err := writer.WriteString(cmd + "\r\n"); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	if err := writer.Flush(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	_, err := readingResp(reader)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func readingResp(reader *bufio.Reader) (string, error) {
	var response strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		response.WriteString(line)
		fmt.Printf("<- %s", line)

		if len(line) < 4 {
			continue
		}

		code := line[:3]
		if line[3] == ' ' {
			if code >= "400" {
				return response.String(), fmt.Errorf("%s", strings.TrimSpace(line))
			}
			break
		}
	}

	return response.String(), nil
}