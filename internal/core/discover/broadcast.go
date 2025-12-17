package discover

import (
	"fmt"
	"log"
	"net"
	"time"
)

func broadcastMe(selfID, selfPort string){
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	for{
		msg := fmt.Sprintf("PEER:ID=%s;PORT=%s", selfID, selfPort)
		conn.Write([]byte(msg))
		time.Sleep(discoverTime)
	}
}