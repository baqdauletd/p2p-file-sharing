package grpcapi

import (
    "context"

    "p2p-file-sharing/internal/grpcapi/generated"
    "p2p-file-sharing/internal/core/peers"
)

type PeerServiceImpl struct {
    generated.UnimplementedPeerServiceServer
    peerManager *peers.PeerManager
}

func (s *PeerServiceImpl) GetCatalog(ctx context.Context, req *generated.PeerID) (*generated.CatalogResponse, error) {
    // files, err := catalog.BuildFileCatalog("shared")
    files, err := s.peerManager.GetCatalog("shared")
    if err != nil {
        return nil, err
    }

    resp := &generated.CatalogResponse{}
    for _, f := range files {
        resp.Files = append(resp.Files, &generated.FileMeta{
            Name: f.Name,
            Size: f.Size,
        })
    }

    return resp, nil
}

func (s *PeerServiceImpl) RequestFile(ctx context.Context, req *generated.File) (*generated.FileResponse, error) {
    // go transfer.StartTCPTransfer(req.FileName, 9000) // non-blocking
    err := s.peerManager.ManageRequestFile(req.FileName, 9000)
    if err != nil {
		return nil, err
	}

    return &generated.FileResponse{Message: "Ready for TCP transfer"}, nil
}

func (s *PeerServiceImpl) GetPeers(ctx context.Context, p *generated.PeerInfo) (*generated.PeersResponse, error) {
    // peers := discover.GetPeersByFileName(f.Name)
    peers := s.peerManager.ManageGetPeers(p.Port)

    resp := &generated.PeersResponse{}
    for _, p := range peers {
        resp.Peers = append(resp.Peers, &generated.PeerInfo{
            Id:   p.ID,
            Ip:   p.IP,
            Port: p.GRPCPort,
        })
    }

    return resp, nil
}
