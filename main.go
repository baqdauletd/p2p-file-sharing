package main

import (
	"flag"
	"fmt"
	"os"
	"p2p-file-sharing/internal/server"
	"p2p-file-sharing/internal/client"
	"p2p-file-sharing/internal/peer"
	"time"
)

func main(){
	mode := flag.String("mode", "", "server or connect")
	port := flag.String("port", "", "port to listen to or connect to")
	host := flag.String("host", "localhost", "host to connect to (only for client)")
	id := flag.String("id", "peer", "IDs for peers")
	flag.Parse()

	switch *mode{
	case "serve":
		fmt.Println("Starting server on port", *port)
		go server.Server(*port)
		go peer.StartDiscovery(*id, *port)

		for {}
	case "connect":
		fmt.Println("Connecting to peer at", *host+":"+*port)
		client.Client(*host, *port)
	case "peers":
		peer.StartDiscovery(*id, *port)
		fmt.Println("Listening for peers... (waiting 6 seconds)")
		time.Sleep(6 * time.Second)
		// fmt.Println("here")
		peerList := peer.GetKnownPeers()
		fmt.Println(len(peerList))
		for _, p := range peerList {
			fmt.Printf("ID: %s, IP: %s, Port: %s\n", p.ID, p.IP, p.Port)
		}
	default:
		fmt.Println("Usage:")
		fmt.Println("  ./p2p-file-sharing -mode=serve -port=8080 -id=peerX")
		fmt.Println("  ./p2p-file-sharing -mode=connect -host=127.0.0.1 -port=8080 -id=peerY")
		fmt.Println("  ./p2p-file-sharing -mode=peers")
		os.Exit(1)
	}
}
