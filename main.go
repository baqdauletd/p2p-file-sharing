package main

import (
	"p2p-file-sharing/internal/core/client"
	"p2p-file-sharing/internal/core/discover"
	"p2p-file-sharing/internal/grpcapi"
)

func main(){
	go discover.Listen(*port)
	go discover.StartDiscovery(*id, *port)
	go client.Client(*host, *port)
	go grpcapi.StartGRPCServer()
}
