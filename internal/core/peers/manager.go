package peers

import (
	"log"
	"sync"
	"time"
)

type Peer struct {
	ID        string
	IP        string
	GRPCPort  string
	LastSeen  time.Time
}

// PeerManager is the single source of truth about peers in the network.
// It does NOT do gRPC calls itself.
type PeerManager struct {
	selfID string

	mu    sync.RWMutex
	peers map[string]*Peer

	connector Connector // interface, not concrete type
	stopCh    chan struct{}
}

// Connector defines what the manager needs from the connector.
// This avoids tight coupling.
type Connector interface {
	ConnectToPeer(peer *Peer)
	DisconnectPeer(peerID string)
}

// NewPeerManager creates a manager instance.
func NewPeerManager(selfID string, connector Connector) *PeerManager {
	return &PeerManager{
		selfID:   selfID,
		peers:   make(map[string]*Peer),
		connector: connector,
		stopCh:  make(chan struct{}),
	}
}

// Start launches all peer-related background loops.
func (pm *PeerManager) Start() {
	go pm.manageListen()
	go pm.manageBroadcast()
	go pm.peerCleanupLoop()
}

// Stop gracefully stops all loops.
func (pm *PeerManager) Stop() {
	close(pm.stopCh)
}

func (pm *PeerManager) ManageGetPeers(port string) ([]Peer){
	peers := pm.GetPeers()
}

// Called by discovery listener when a peer announcement is received
func (pm *PeerManager) HandleDiscoveredPeer(id, ip, grpcPort string) {
	if id == pm.selfID {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	peer, exists := pm.peers[id]
	if exists {
		peer.LastSeen = time.Now()
		return
	}

	peer = &Peer{
		ID:       id,
		IP:       ip,
		GRPCPort: grpcPort,
		LastSeen: time.Now(),
	}

	pm.peers[id] = peer
	log.Println("New peer discovered:", id)

	// Delegate connection responsibility
	pm.connector.ConnectToPeer(peer)
}

func (pm *PeerManager) ListPeers() []*Peer {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	out := make([]*Peer, 0, len(pm.peers))
	for _, p := range pm.peers {
		out = append(out, p)
	}
	return out
}

func (pm *PeerManager) GetPeer(id string) (*Peer, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	p, ok := pm.peers[id]
	return p, ok
}

func (pm *PeerManager) peerCleanupLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.cleanupStalePeers()
		case <-pm.stopCh:
			return
		}
	}
}

func (pm *PeerManager) cleanupStalePeers() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	now := time.Now()
	for id, peer := range pm.peers {
		if now.Sub(peer.LastSeen) > 30*time.Second {
			log.Println("Peer timed out:", id)
			delete(pm.peers, id)
			pm.connector.DisconnectPeer(id)
		}
	}
}

func (pm *PeerManager) manageListen() {
	for {
		select {
		case <-pm.stopCh:
			return
		default:
			// Example:
			// id, ip, grpcPort := discovery.Receive()
			// pm.HandleDiscoveredPeer(id, ip, grpcPort)
		}
	}
}

func (pm *PeerManager) manageBroadcast() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Example:
			// discovery.Broadcast(pm.selfID, grpcPort)
		case <-pm.stopCh:
			return
		}
	}
}




