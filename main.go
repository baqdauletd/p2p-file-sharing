package main

import (
	"flag"
	"fmt"
	"os"
	"p2p-file-sharing/internal/client"
	"p2p-file-sharing/internal/server"
)

func main(){
	mode := flag.String("mode", "", "server or connect")
	port := flag.String("port", "", "port to listen to or connect to")
	host := flag.String("host", "localhost", "host to connect to (only for client)")
	flag.Parse()

	switch *mode{
	case "serve":
		fmt.Println("Starting server on port", *port)
		server.Server(*port)
	case "connect":
		fmt.Println("Connecting to peer at", *host+":"+*port)
		client.Client(*host, *port)
	default:
		fmt.Println("Usage:")
		fmt.Println("  ./p2p-file-sharing -mode=serve -port=8080")
		fmt.Println("  ./p2p-file-sharing -mode=connect -host=127.0.0.1 -port=8080")
		os.Exit(1)
	}
}
