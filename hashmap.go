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

const MaxInt = int(^uint(0) >> 1)

type HashMap struct {
	sync.RWMutex

	size        int64
	nodes       []*Node
	loadFactor  float64
	resizeTimes int
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
	prev *Entry
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
	m.resize(0)
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

	//If key does not exist
	n.Lock()
	if m.setNodeEntry(n, &Entry{k: k, p: unsafe.Pointer(&v), hash: h}) {
		n.size++
		atomic.AddInt64(&m.size, 1)
	}
	n.Unlock()
	return nil
}

func (m *HashMap) MSet(ks []interface{}, vs []interface{}){
	if len(ks) != len(vs) {
		return
	}
	m.resize(int(float64(len(ks)) * 0.66))
	for i, k := range ks {
		m.Set(k, vs[i])
	}
}

func (m *HashMap) setNodeEntry(n *Node, e *Entry) bool {
	if n.head == nil {
		n.head = e
		n.tail = e
	} else {
		next := n.head
		for next != nil {
			if next.k == e.k {
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

func (m *HashMap) dilate(c int) bool {
	return (m.size + int64(c)) > int64(float64(len(m.nodes))*m.loadFactor*3) && len(m.nodes) * m.multiple(c) <= MaxInt
}

func (m *HashMap) multiple(c int) int {
	expect := int64(float64(m.size + int64(c)) / 3 / m.loadFactor)
	l := len(m.nodes)
	mul := 2
	for int64(l * mul) < expect {
		mul <<= 1
	}
	return mul
}

func (m *HashMap) resize(c int) {
	if m.dilate(c){
		m.Lock()
		for m.dilate(c){
			m.doResize(m.multiple(c))
		}
		m.Unlock()
	}
}

func (m *HashMap) doResize(multiple int) {
	capacity := len(m.nodes) * multiple
	nodes := allocate(capacity)
	var size int64 = 0
	for _, old := range m.nodes {
		next := old.head
		for next != nil {
			newNode := nodes[indexOf(next.hash, capacity)]
			e := next.clone()
			if newNode.head == nil {
				newNode.head = e
				newNode.tail = e
			} else {
				newNode.tail.next = e
				e.prev = newNode.tail
				newNode.tail = e
			}
			size++
			newNode.size++
			next = next.next
		}
	}
	m.nodes = nodes
	m.size = size
	m.resizeTimes ++
}

func (m *HashMap) getNodeEntry(n *Node, k interface{}) *Entry {
	if n != nil {
		next := n.head
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
	e := m.getNodeEntry(n, k)
	if e != nil {
		if e.prev == nil && e.next == nil {
			n.head = nil
			n.tail = nil
		} else if e.prev == nil {
			n.head = e.next
			e.next.prev = nil
		} else if e.next == nil {
			n.tail = e.prev
			e.prev.next = nil
		} else {
			e.prev.next = e.next
			e.next.prev = e.prev
		}
		n.size--
		atomic.AddInt64(&m.size, -1)
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
