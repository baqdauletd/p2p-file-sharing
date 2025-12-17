package grpcapi

import (
    "log"
    "net"

    "google.golang.org/grpc"
    "p2p-file-sharing/internal/grpcapi/generated"
)

func StartGRPCServer(port string) {
    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()

    generated.RegisterPeerServiceServer(grpcServer, &PeerServiceImpl{})

    log.Println("Peer and Tracker gRPC servers running on port", port)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}