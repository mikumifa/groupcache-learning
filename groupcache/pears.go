package groupcache

import pb "geecache-learning/groupcache/cachepb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
// 选择Peer
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// 从Peer中选择缓存
type PeerGetter interface {
	//需要group和key
	Get(in *pb.Request, out *pb.Response) error
}
