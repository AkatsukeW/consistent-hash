package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Consistent struct {
	hasSortedNode []uint32
	circle        map[uint32]string
	nodes         map[string]struct{}
	sync.RWMutex
	// count virtual nodes
	virtualNodeCount int
}

func (c *Consistent) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) AddNode(node string, virtualNodeCount int) error {
	if node == "" {
		return nil
	}
	c.Lock()
	defer c.Unlock()

	if c.circle == nil {
		c.circle = make(map[uint32]string)
	}

	if c.nodes == nil {
		c.nodes = make(map[string]struct{})
	}
	c.nodes[node] = struct{}{}

	// add virtual node
	for i := 0; i < virtualNodeCount; i++ {
		virtualKey := c.hashKey(node + strconv.Itoa(i))
		c.circle[virtualKey] = node
		c.hasSortedNode = append(c.hasSortedNode, virtualKey)
	}

	// sort
	sort.Slice(c.hasSortedNode, func(i, j int) bool {
		return c.hasSortedNode[i] < c.hasSortedNode[j]
	})
	return nil
}

func (c *Consistent) GetNode(key string) string {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashKey(key)
	index := c.getPosition(hash)

	return c.circle[c.hasSortedNode[index]]
}

func (c *Consistent) getPosition(hash uint32) int {
	index := sort.Search(len(c.hasSortedNode), func(i int) bool {
		return c.hasSortedNode[i] >= hash
	})

	if index < len(c.hasSortedNode) {
		if index == len(c.hasSortedNode) {
			return 0
		}

		return index
	}
	return len(c.hasSortedNode) - 1
}
