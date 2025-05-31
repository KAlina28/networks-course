package main

import (
	"encoding/binary"
	"lab08/internal"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	LossPercentage = 30
)

func main() {
	rand.Seed(time.Now().UnixNano())

	conn, err := net.ListenPacket("udp", ":8888")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	defer conn.Close()

	var (
		receivedData   []byte
		expectedSeqNum uint32
		lastAddr       net.Addr
		lastActive     = time.Now()
	)

	for {
		if time.Since(lastActive) > 30*time.Second {
			break
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		buf := make([]byte, 2048)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Printf("Read error: %v", err)
			continue
		}

		lastActive = time.Now()
		lastAddr = addr

		if rand.Intn(100) < LossPercentage {
			log.Println("packet lost from client")
			continue
		}

		pkt, err := internal.ParsePacket(buf[:n])
		if err != nil || !internal.VerifyChecksum(pkt) {
			log.Println("corrupt packet from client")
			continue
		}

		if pkt.SeqNum == expectedSeqNum {
			receivedData = append(receivedData, pkt.Data...)
			expectedSeqNum = 1 - expectedSeqNum
		}

		ack := internal.MakeACK(pkt.SeqNum)
		if rand.Intn(100) >= LossPercentage {
			conn.WriteTo(ack, addr)
		}
	}

	if len(receivedData) > 0 {
		err := os.WriteFile("received_file.txt", receivedData, 0644)
		if err != nil {
			log.Printf("error saving file: %v", err)
		} else {
			log.Printf("saved to %s", "received_file.txt")
		}
	}

	if lastAddr != nil {
		log.Println("sending file back to client.")
		err := sendFileToClient(conn, lastAddr, "server_response.txt")
		if err != nil {
			log.Printf("Send error: %v", err)
		}
	}
}

func sendFileToClient(conn net.PacketConn, addr net.Addr, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	packets := splitIntoPackets(data)
	for i, pkt := range packets {
		for {
			if rand.Intn(100) >= LossPercentage {
				_, err := conn.WriteTo(pkt, addr)
				if err != nil {
					continue
				}
				log.Printf("sent packet %d", i)
			} else {
				log.Printf("packet %d lost (simulated)", i)
			}

			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			buf := make([]byte, 8)
			_, _, err := conn.ReadFrom(buf)
			if err != nil {
				continue
			}

			if rand.Intn(100) < LossPercentage {
				log.Printf("ACK lost for packet %d", i)
				continue
			}

			ackSeq := binary.BigEndian.Uint32(buf[:4])
			if ackSeq == uint32(i%2) {
				break
			}
		}
	}
	return nil
}

func splitIntoPackets(data []byte) [][]byte {
	var packets [][]byte
	seqNum := 0

	for offset := 0; offset < len(data); offset += 1024 {
		end := offset + 1024
		if end > len(data) {
			end = len(data)
		}
		chunk := data[offset:end]
		pkt := internal.MakePacket(uint32(seqNum%2), chunk)
		packets = append(packets, pkt)
		seqNum++
	}
	return packets
}
