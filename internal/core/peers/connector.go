package peers

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "p2p-file-sharing/internal/grpcapi/generated"
)

type GRPCConnector struct {
	mu sync.Mutex

	clients map[string]pb.PeerServiceClient
	conns   map[string]*grpc.ClientConn
}


func NewGRPCConnector() *GRPCConnector {
	return &GRPCConnector{
		clients: make(map[string]pb.PeerServiceClient),
		conns:   make(map[string]*grpc.ClientConn),
	}
}

func (c *GRPCConnector) ConnectToPeer(peer *Peer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.clients[peer.ID]; exists {
		return // already connected
	}

	addr := peer.IP + ":" + peer.GRPCPort

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Println("Failed to connect to peer:", peer.ID, err)
		return
	}

	client := pb.NewPeerServiceClient(conn)

	c.conns[peer.ID] = conn
	c.clients[peer.ID] = client

	log.Println("Connected to peer via gRPC:", peer.ID)
}

func (c *GRPCConnector) DisconnectPeer(peerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, ok := c.conns[peerID]
	if !ok {
		return
	}

	conn.Close()
	delete(c.conns, peerID)
	delete(c.clients, peerID)

	log.Println("Disconnected from peer:", peerID)
}

func (c *GRPCConnector) GetClient(peerID string) (pb.PeerServiceClient, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := c.clients[peerID]
	return client, ok
}
