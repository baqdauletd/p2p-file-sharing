package peer

type Peer struct {
    ID   string
    IP   string
    Port string
}

var knownPeers = make(map[string]Peer)