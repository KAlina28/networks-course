package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jlaffaye/ftp"
)

func main() {
	client, err := ftp.Dial(
		"ftp.dlptest.com:21",
		ftp.DialWithDisabledEPSV(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	err = client.Login("dlpuser", "rNrKYTX9g7z3RgJRmxWuGHbeu")
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("1 - список всего")
		fmt.Println("2 - загрузка")
		fmt.Println("3 - выгрузка")
		fmt.Println("4 - конец")

		var choice int
		fmt.Print("-> ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			log.Println(err)
			continue
		}

		switch choice {
		case 1:
			listFiles(client)
		case 2:
			uploadFile(client)
		case 3:
			downloadFile(client)
		case 4:
			fmt.Println("пока")
			return
		default:
			fmt.Println("нет такого")
		}
	}
}

func listFiles(client *ftp.ServerConn) {
	listFiles, err := client.List("")
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range listFiles {
		fmt.Printf("-  %s\n", file.Name)
	}
}

func uploadFile(client *ftp.ServerConn) {
	var filePath string
	fmt.Print("введи путь: ")

	_, err := fmt.Scanln(&filePath)
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)

	err = client.Stor(fileName, file)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("загружено %s \n", fileName)
}

func downloadFile(client *ftp.ServerConn) {
	var fileName string
	fmt.Print("введи название файла:")
	fmt.Scan(&fileName)

	reader, err := client.Retr(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer reader.Close()

	localFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, reader)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("скачан %s \n", fileName)
}
