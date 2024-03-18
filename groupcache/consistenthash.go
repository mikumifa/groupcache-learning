package groupcache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// peersMap constains all hashed keys
type peersMap struct {
	hash     Hash  // hash函数
	replicas int   //虚拟节点的倍数， 应该真实节点有多少个虚拟节点
	keys     []int // Sorted
	hashMap  map[int]string
}

// New creates a peersMap instance
// 可以自定义Hash
func New(replicas int, fn Hash) *peersMap {
	m := &peersMap{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
// 添加key到hash里面，同时存储
// 重叠也无所谓，反正是以hashMap的为主
// 添加一个节点的名字
func (m *peersMap) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
// 这个相当于告诉吴key是string时候应该选哪个节点名字
func (m *peersMap) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//找不到就是n
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
