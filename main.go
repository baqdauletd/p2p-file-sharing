package main

import (
	"flag"
	"fmt"
	"os"
	"p2p-file-sharing/internal/pp"
	"time"
)

func main(){
	mode := flag.String("mode", "", "server or connect")
	port := flag.String("port", "", "port to listen to or connect to")
	host := flag.String("host", "localhost", "host to connect to (only for client)")
	id := flag.String("id", "peer", "IDs for peers")
	udpPort := flag.String("udpport", "9999", "UDP port for peer discovery")
	flag.Parse()

	switch *mode{
	case "serve":
		fmt.Println("Starting server on port", *port)
		go pp.Server(*port)
		go pp.StartDiscovery(*id, *port, *udpPort)

		for {}
	case "connect":
		fmt.Println("Connecting to peer at", *host+":"+*port)
		pp.Client(*host, *port)
	case "peers":
		pp.StartDiscovery(*id, *port, *udpPort)
		fmt.Println("Listening for peers... (waiting 6 seconds)")
		time.Sleep(6 * time.Second)
		// fmt.Println("here")
		peerList := pp.GetKnownPeers()
		fmt.Println(len(peerList))
		for _, p := range peerList {
			// fmt.Println("here222")
			fmt.Printf("ID: %s, IP: %s, Port: %s\n", p.ID, p.IP, p.Port)
		}
	default:
		fmt.Println("Usage:")
		fmt.Println("  ./p2p-file-sharing -mode=serve -port=8080")
		fmt.Println("  ./p2p-file-sharing -mode=connect -host=127.0.0.1 -port=8080")
		os.Exit(1)
	}
}
