package hashmap

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type HashMap struct {
	sync.RWMutex

	size       int64
	nodes      []*Node
	loadFactor float64
}

type Node struct {
	sync.Mutex

	head *Entry
	tail *Entry
	size int
}

type Entry struct {
	k    interface{}
	p    unsafe.Pointer
	hash uint64
	next *Entry
}

func New() *HashMap {
	return &HashMap{
		nodes:      allocate(16),
		loadFactor: 0.7,
	}
}

func allocate(capacity int) (nodes []*Node) {
	nodes = make([]*Node, capacity)
	for i := 0; i < capacity; i++ {
		nodes[i] = &Node{}
	}
	return
}

func hash(k interface{}) uint64 {
	if k == nil {
		return 0
	}
	switch x := k.(type) {
	case string:
		return bytesHash([]byte(x))
	case []byte:
		return bytesHash(x)
	case bool:
		if x {
			return 0
		} else {
			return 1
		}
	case time.Time:
		return uint64(x.UnixNano())
	case int:
		return uint64(x)
	case int8:
		return uint64(x)
	case int16:
		return uint64(x)
	case int32:
		return uint64(x)
	case int64:
		return uint64(x)
	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case float32:
		return math.Float64bits(float64(x))
	case float64:
		return math.Float64bits(x)
	case uintptr:
		return uint64(x)
	}
	panic("unsupported key type.")
}

func bytesHash(bytes []byte) uint64 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(bytes)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(bytes[i])
	}
	return uint64(hash)
}

func indexOf(hash uint64, capacity int) int {
	return int(hash & uint64(capacity-1))
}

//Set will CAS the existing value if k exists. If k is new, this function is locked and set node's head
//Similar to Java's hashmap's Put
//returns old value if k previously exists
//returns nil if k is new
func (m *HashMap) Set(k interface{}, v interface{}) interface{} {
	m.resize()
	m.RLock()
	defer m.RUnlock()

	h, nodes := hash(k), m.nodes
	n := nodes[indexOf(h, len(nodes))]

	//If key exists
	if e := m.getNodeEntry(n, k); e != nil {
		oldValue := e.value()
		atomic.StorePointer(&e.p, unsafe.Pointer(&v))
		return oldValue
	}
	n.Lock()
	if m.setNodeEntry(n, &Entry{k: k, p: unsafe.Pointer(&v), hash: h}) {
		n.size++
		atomic.AddInt64(&m.size, 1)
	}
	n.Unlock()
	return nil
}

func (m *HashMap) setNodeEntry(n *Node, e *Entry) bool {
	if n.head == nil {
		n.head, n.tail = e, e
	} else {
		next := n.head
		for next != nil {
			if next.k == e.k {
				next.p = e.p
				return false
			}
			next = next.next
		}
		n.tail.next, n.tail = e, e
	}
	return true
}

func (m *HashMap) dilate() bool {
	return m.size > int64(float64(len(m.nodes))*m.loadFactor*3)
}

func (m *HashMap) resize() {
	if m.dilate() {
		m.Lock()
		defer m.Unlock()
		if m.dilate() {
			m.doResize()
		}
	}
}

func (m *HashMap) doResize() {
	capacity := len(m.nodes) * 2
	nodes := allocate(capacity)
	size := int64(0)
	for _, old := range m.nodes {
		next := old.head
		for next != nil {
			newNode, e := nodes[indexOf(next.hash, capacity)], next.clone()
			if newNode.head == nil {
				newNode.head, newNode.tail = e, e
			} else {
				newNode.tail.next, newNode.tail = e, e
			}
			size++
			newNode.size++
			next = next.next
		}
	}
	m.nodes = nodes
	m.size = size
}

func (m *HashMap) getNodeEntry(n *Node, k interface{}) *Entry {
	next := n.head
	for next != nil {
		if next.k == k {
			return next
		}
		next = next.next
	}
	return nil
}

func (m *HashMap) Get(k interface{}) (interface{}, bool) {
	nodes := m.nodes
	n := nodes[indexOf(hash(k), len(nodes))]
	if n != nil {
		e := m.getNodeEntry(n, k)
		if e != nil {
			return e.value(), true
		}
	}
	return nil, false
}

func (m *HashMap) Del(k interface{}) bool {
	m.RLock()
	defer m.RUnlock()

	nodes := m.nodes
	n := nodes[indexOf(hash(k), len(nodes))]
	n.Lock()
	defer n.Unlock()
	var next, prev *Entry = n.head, nil
	for next != nil {
		if next.k == k {
			if prev == nil {
				n.head = next.next
				if n.head == nil {
					n.tail = nil
				}
			}else{
				prev.next = next.next
				if prev.next == nil {
					n.tail = prev
				}
			}
			n.size --
			atomic.AddInt64(&m.size, -1)
			return true
		}
		prev, next = next, next.next
	}
	return false
}

func (e *Entry) clone() *Entry {
	return &Entry{
		k:    e.k,
		p:    e.p,
		hash: e.hash,
	}
}

func (e *Entry) value() interface{} {
	return *(*interface{})(e.p)
}

func (m *HashMap) UnmarshalJSON(b []byte) error {
	data := map[string]interface{}{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	for k, v := range data {
		m.Set(k, v)
	}
	return nil
}

func (m *HashMap) MarshalJSON() ([]byte, error) {
	nodes := m.nodes
	data := map[string]interface{}{}
	for _, node := range nodes {
		next := node.head
		for next != nil {
			data[fmt.Sprintf("%v", next.k)] = next.value()
			next = next.next
		}
	}
	return json.Marshal(data)
}
