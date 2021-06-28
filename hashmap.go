package _struct

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type HashMap struct {
	sync.RWMutex

	size int64
	nodes []*Node
	loadFactor float64
}

type Node struct {
	sync.Mutex

	header *Entry
	tail *Entry
	size int
}

type Entry struct {
	k interface{}
	p unsafe.Pointer
	hash int
	next *Entry
	prev *Entry
}

func New() *HashMap{
	return &HashMap{
		nodes: allocate(16),
		loadFactor: 0.7,
	}
}

func allocate(capacity int) (nodes []*Node){
	nodes = make([]*Node, capacity)
	for i := 0; i < capacity; i ++{
		nodes[i] = &Node{}
	}
	return
}

func hash(k interface{}) int {
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
		}else{
			return 1
		}
	case time.Time:
		return x.Nanosecond()
	case int:
		return int(uintptr(x))
	case int8:
		return int(uintptr(x))
	case int16:
		return int(uintptr(x))
	case int32:
		return int(uintptr(x))
	case int64:
		return int(uintptr(x))
	case uint:
		return int(uintptr(x))
	case uint8:
		return int(uintptr(x))
	case uint16:
		return int(uintptr(x))
	case uint32:
		return int(uintptr(x))
	case uint64:
		return int(uintptr(x))
	case uintptr:
		return int(x)
	}
	panic("unsupported key type.")
}

func bytesHash(bytes []byte) int{
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(bytes)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(bytes[i])
	}
	return int(hash)
}

func indexOf(hash int, capacity int) int{
	return hash & (capacity - 1)
}

func (m *HashMap) Set(k interface{}, v interface{}) interface{} {
	m.resize()
	m.RLock()
	defer m.RUnlock()

	h, nodes := hash(k), m.nodes
	n := nodes[indexOf(h, len(nodes))]

	if e := m.getNodeEntry(n, k); e != nil {
		for {
			p := atomic.LoadPointer(&e.p)
			if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(&v)) {
				return v
			}
		}
	}

	n.Lock()
	defer n.Unlock()
	if m.setNodeEntry(n, &Entry{k: k, p: unsafe.Pointer(&v), hash: h}) {
		n.size ++
		atomic.AddInt64(&m.size, 1)
	}
	return v
}

func (m *HashMap) setNodeEntry(n *Node, e *Entry) bool{
	if n.header == nil {
		n.header = e
		n.tail = e
	} else {
		next := n.header
		for next != nil{
			if next.k == e.k{
				next.p = e.p
				return false
			}
			next = next.next
		}
		n.tail.next = e
		e.prev = n.tail
		n.tail = e
	}
	return true
}

func (m *HashMap) dilate() bool {
	return m.size > int64(float64(len(m.nodes)) * m.loadFactor * 3)
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

func (m *HashMap) doResize()  {
	capacity := len(m.nodes) * 2
	nodes := allocate(capacity)
	size := int64(0)
	for _, old := range m.nodes {
		next := old.header
		for next != nil {
			newNode := nodes[indexOf(next.hash, capacity)]
			e := next.clone()
			if newNode.header == nil {
				newNode.header = e
				newNode.tail = e
			}else{
				newNode.tail.next = e
				e.prev = newNode.tail
				newNode.tail = e
			}
			size ++
			newNode.size ++
			next = next.next
		}
	}
	m.nodes = nodes
	m.size = size
}

func (m *HashMap) getNodeEntry(n *Node, k interface{}) *Entry {
	if n != nil {
		next := n.header
		for next != nil {
			if next.k == k {
				return next
			}
			next = next.next
		}
	}
	return nil
}

func (m *HashMap) Get(k interface{}) (interface{}, bool) {
	nodes := m.nodes
	n := nodes[indexOf(hash(k), len(nodes))]
	if n != nil {
		e := m.getNodeEntry(n, k)
		if e != nil {
			return *(*interface{})(e.p), true
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
	e := m.getNodeEntry(n, k)
	if e != nil {
		if e.prev == nil && e.next == nil{
			n.header = nil
			n.tail = nil
		}else if e.prev == nil {
			n.header = e.next
			e.next.prev = nil
		}else if e.next == nil {
			n.tail = e.prev
			e.prev.next = nil
		}else{
			e.prev.next = e.next
			e.next.prev = e.prev
		}
		n.size --
		atomic.AddInt64(&m.size, -1)
	}
	return false
}

func (e *Entry) clone() *Entry{
	return &Entry{
		k: e.k,
		p: e.p,
		hash: e.hash,
	}
}