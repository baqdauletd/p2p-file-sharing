package pp

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
	broadcastPort = 9999
	broadcastAddr = "255.255.255.255:9999"
	discoverTime = 5*time.Second
)


type PeerInfo struct {
	ID   string
	IP   string
	Port string
}

var knownPeers = make(map[string]PeerInfo)

func StartDiscovery(selfID, selfPort, selfUDPport string) {
	// log.Println("here1")
	go listenForPeers(selfUDPport)
	go broadcastPresence(selfID, selfPort)
}

func broadcastPresence(selfID, selfPort string){
	// log.Println("here")
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

func listenForPeers(selfUDPport string){
	// addr, err := net.ResolveUDPAddr("udp", ":"+selfUDPport)
	// if err != nil {
	// 	log.Fatalf("Failed to resolve UDP address: %v", err)
	// }
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	log.Fatalf("Failed to Listen UDP: %v", err)
	// }
	conn, err := listenUDPWithReuse(selfUDPport)
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
			peerInfo := parsePeerInfo(data, remoteAddr)
			if peerInfo.ID != "" && peerInfo.ID != "self" {
				if _, exists := knownPeers[peerInfo.ID]; !exists {
					knownPeers[peerInfo.ID] = peerInfo
					log.Println("Inserted peer:", peerInfo.ID)
				}
			}
			// for _, p := range knownPeers{
			// 	fmt.Printf("ID: %s, IP: %s, Port: %s\n", p.ID, p.IP, p.Port)
			// }
			// fmt.Println()
		}
	}
}

func parsePeerInfo(data string, remoteAddr *net.UDPAddr) PeerInfo{
	parts := strings.Split(data, ";")
	id := strings.TrimPrefix(parts[0], "PEER:ID=")
	port := strings.TrimPrefix(parts[1], "PORT=")

	return PeerInfo{
		ID: id,
		IP: remoteAddr.String(),
		Port: port,
	}
}

func GetKnownPeers() []PeerInfo{
	// fmt.Println("here333")
	list := []PeerInfo{}
	for _, peerInfo := range knownPeers{
		// fmt.Println("her444")
		list = append(list, peerInfo)
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

	// Cast to *net.UDPConn
	conn, ok := pc.(*net.UDPConn)
	if !ok {
		pc.Close()
		return nil, fmt.Errorf("failed to cast to UDPConn")
	}
	return conn, nil
}

