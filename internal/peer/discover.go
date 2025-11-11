package peer

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"
	"time"
)

const(
	udpAddr = "9999"
	broadcastAddr = "255.255.255.255:9999"
	discoverTime = 5*time.Second
)


// type Peer struct {
// 	ID   string
// 	IP   string
// 	Port string
// }

func StartDiscovery(selfID, selfPort string) {
	go listenForPeers()
	go broadcastPresence(selfID, selfPort)
}

func broadcastPresence(selfID, selfPort string){
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	// log.Println("port:", selfPort)
	// log.Println("id:", selfID)

	for{
		msg := fmt.Sprintf("PEER:ID=%s;PORT=%s", selfID, selfPort)
		conn.Write([]byte(msg))
		time.Sleep(discoverTime)
	}
}

func listenForPeers(){
	// addr, err := net.ResolveUDPAddr("udp", ":"+selfUDPport)
	// if err != nil {
	// 	log.Fatalf("Failed to resolve UDP address: %v", err)
	// }
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	log.Fatalf("Failed to Listen UDP: %v", err)
	// }
	conn, err := listenUDPWithReuse(udpAddr)
	if err != nil {
		log.Fatalf("Failed to Listen UDP: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		data := string(buf[:n])
		if strings.HasPrefix(data, "PEER:") {
			Peer := parsePeer(data, remoteAddr)
			if Peer.ID != "" && Peer.ID != "self" {
				if _, exists := knownPeers[Peer.ID]; !exists {
					knownPeers[Peer.ID] = Peer
					fmt.Println("Inserted peer:", Peer.ID)
				}
			}
			// for _, p := range knownPeers{
			// 	fmt.Printf("ID: %s, IP: %s, Port: %s\n", p.ID, p.IP, p.Port)
			// }
			// fmt.Println()
		}
	}
}

func parsePeer(data string, remoteAddr *net.UDPAddr) Peer{
	parts := strings.Split(data, ";")
	id := strings.TrimPrefix(parts[0], "PEER:ID=")
	port := strings.TrimPrefix(parts[1], "PORT=")

	return Peer{
		ID: id,
		IP: remoteAddr.String(),
		Port: port,
	}
}

func GetKnownPeers() []Peer{
	// fmt.Println("here333")
	list := []Peer{}
	for _, Peer := range knownPeers{
		// fmt.Println("her444")
		list = append(list, Peer)
	}

	return list
}



// for using the same UDPaddress by several terminals, since net.ListenUDP doesn't allow to do so
func listenUDPWithReuse(port string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return nil, err
	}

	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			c.Control(func(fd uintptr) {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			})
			return err
		},
	}

	pc, err := lc.ListenPacket(context.Background(), "udp", addr.String())
	if err != nil {
		return nil, err
	}

	// cast to *net.UDPConn
	conn, ok := pc.(*net.UDPConn)
	if !ok {
		pc.Close()
		return nil, fmt.Errorf("failed to cast to UDPConn")
	}
	return conn, nil
}

