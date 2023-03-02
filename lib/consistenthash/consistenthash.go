package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int
	replicates  int // 虚拟节点数
	nodehashMap map[int]string
}

func NewNodeMap(fn HashFunc, replicates int) *NodeMap {
	m := &NodeMap{
		hashFunc:    fn,
		replicates:  replicates,
		nodehashMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashs) == 0
}

func (m *NodeMap) AddNodes(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		// 为每个node添加多个虚拟节点
		for i := 0; i < m.replicates; i++ {
			hash := int(m.hashFunc([]byte(key + strconv.Itoa(i))))
			m.nodeHashs = append(m.nodeHashs, hash)
			m.nodehashMap[hash] = key
		}
	}
	sort.Ints(m.nodeHashs)
}

func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	i := sort.Search(len(m.nodeHashs), func(i int) bool {
		return m.nodeHashs[i] >= hash
	})
	return m.nodehashMap[m.nodeHashs[i%len(m.nodeHashs)]]
}
