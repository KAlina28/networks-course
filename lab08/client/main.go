package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
	"lab08/internal"
)

const (
	ServerAddr     = "127.0.0.1:8888"
	LossPercentage = 30
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not cprrect input")
		return
	}

	data, _ := os.ReadFile(os.Args[1])

	rand.Seed(time.Now().UnixNano())

	conn, _ := net.Dial("udp", ServerAddr)
	defer conn.Close()

	seq := uint32(0)
	for offset := 0; offset < len(data); {
		end := min(offset + 1024, len(data))

		packet := internal.MakePacket(seq, data[offset:end])

		for {
			if rand.Intn(100) >= LossPercentage {
				_, err := conn.Write(packet)
				if err != nil {
					log.Printf("send error: %v", err)
					continue
				}
				log.Printf("sent packet %d", seq)
			} else {
				log.Printf("imulated loss of packet %d", seq)
			}

			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			ackBuf := make([]byte, 8)
			_, err := conn.Read(ackBuf)
			if err != nil {
				log.Printf("timeout/wrong ACK for packet %d", seq)
				continue
			}

			ackSeq, ackChecksum := binary.BigEndian.Uint32(ackBuf[:4]), binary.BigEndian.Uint32(ackBuf[4:])
			if ackSeq == seq && ackChecksum == internal.CalculateChecksum(ackSeq, nil) {
				log.Printf("received ACK for packet %d", seq)
				break
			} else {
				log.Printf("corrupt ACK or wrong seq: got %d", ackSeq)
			}
		}
		offset += 1024
		seq = 1 - seq
	}

	log.Println("finished sending. waiting for file from server.")

	var receivedData []byte
	expectedSeq := uint32(0)
	timeout := time.Now().Add(35 * time.Second)
	for {
		if time.Now().After(timeout) {
			log.Println("timed out waiting for server packets.")
			break
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}

		if rand.Intn(100) < LossPercentage {
			log.Println("simulated loss of packet from server")
			continue
		}

		pkt, err := internal.ParsePacket(buf[:n])
		if err != nil || !internal.VerifyChecksum(pkt) {
			log.Println("sorrupt packet from server")
			continue
		}

		log.Printf("received packet %d from server", pkt.SeqNum)
		if pkt.SeqNum == expectedSeq {
			receivedData = append(receivedData, pkt.Data...)
			expectedSeq = 1 - expectedSeq
			timeout = time.Now().Add(10 * time.Second)
		}

		ack := internal.MakeACK(pkt.SeqNum)
		if rand.Intn(100) >= LossPercentage {
			conn.Write(ack)
		} else {
			log.Printf("simulated loss of ACK %d to server", pkt.SeqNum)
		}
	}

	if len(receivedData) > 0 {
		err := os.WriteFile("client_received.txt", receivedData, 0644)
		if err != nil {
			log.Printf("failed to save: %v", err)
		} else {
			log.Printf("saved to %s", "client_received.txt")
		}
	}
}
