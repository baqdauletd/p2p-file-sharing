package discover
type Peer struct {
    ID   string
    IP   string
    Port string
}

var knownPeers = make(map[string]Peer)
var peerFiles = make(map[string][]string) // peerID → list of file names

func AddPeer(p Peer) {
    knownPeers[p.ID] = p
}

// RegisterFiles tells tracker which files a peer has
func RegisterFiles(peerID string, files []string) {
    peerFiles[peerID] = files
}

// GetPeersByFileName returns peers who have the given file
func GetPeersByFileName(name string) []Peer {
    result := []Peer{}
    for pid, fileList := range peerFiles {
        for _, f := range fileList {
            if f == name {
                if p, ok := knownPeers[pid]; ok {
                    result = append(result, p)
                }
            }
        }
    }
    return result
}
