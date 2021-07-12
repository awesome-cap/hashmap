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

const MaxInt = 2147483647

type HashMap struct {
	sync.RWMutex

	size       int64
	table      *Table
	loadFactor float64
}

type Table struct {
	nodes []*Node
	ab    int
}

type Node struct {
	sync.Mutex

	head *Entry
	tail *Entry
	size int64
}

type Entry struct {
	k    interface{}
	p    unsafe.Pointer
	hash uint64
	flag int32 // 1 deleted
	next []*Entry
	prev []*Entry
}

func New() *HashMap {
	return &HashMap{
		table: &Table{
			nodes: allocate(16),
			ab:    0,
		},
		loadFactor: 0.7 * 3,
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

func (t *Table) len() int {
	return len(t.nodes)
}

func (m *HashMap) Size() int64 {
	return m.size
}

//Set will CAS the existing value if k exists. If k is new, this function is locked and set node's head
//Similar to Java's hashmap's Put
//returns old value if k previously exists
//returns nil if k is new
func (m *HashMap) Set(k interface{}, v interface{}) interface{} {
	m.resize()
	m.RLock()
	defer m.RUnlock()

	h, t := hash(k), m.table
	n := t.nodes[indexOf(h, t.len())]

	//If key exists
	if e := m.getNodeEntry(t, n, k); e != nil {
		oldValue := e.Value()
		atomic.StorePointer(&e.p, unsafe.Pointer(&v))
		return oldValue
	}
	n.Lock()
	if m.setNodeEntry(t, n, &Entry{k: k, p: unsafe.Pointer(&v), hash: h, next: make([]*Entry, 2), prev: make([]*Entry, 2)}, false) {
		atomic.AddInt64(&n.size, 1)
		atomic.AddInt64(&m.size, 1)
	}
	n.Unlock()
	return nil
}

func (m *HashMap) MSet(ks []interface{}, vs []interface{}) {
	if len(ks) != len(vs) {
		return
	}
	for i, k := range ks {
		m.Set(k, vs[i])
	}
}

func (m *HashMap) SetNX(k interface{}, v interface{}) bool {
	m.resize()
	m.RLock()
	defer m.RUnlock()
	t := m.table
	n, h := t.getKeyNode(k)
	n.Lock()
	defer n.Unlock()
	return m.setNodeEntry(t, n, &Entry{k: k, p: unsafe.Pointer(&v), hash: h, next: make([]*Entry, 2), prev: make([]*Entry, 2)}, true)
}

func (t *Table) getKeyNode(k interface{}) (*Node, uint64) {
	h, nodes := hash(k), t.nodes
	i := indexOf(h, len(nodes))
	return nodes[i], h
}

func (m *HashMap) setNodeEntry(t *Table, n *Node, e *Entry, nx bool) bool {
	if n.head == nil {
		n.head, n.tail = e, e
	} else {
		next := n.head
		for next != nil {
			if next.k == e.k {
				if !nx {
					next.p = e.p
				}
				return false
			}
			next = next.next[t.ab]
		}
		n.tail.next[t.ab], e.prev[t.ab], n.tail = e, n.tail, e
	}
	return true
}

func (m *HashMap) dilate() bool {
	return m.size > int64(float64(m.table.len())*m.loadFactor) && m.table.len()*2 <= MaxInt
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
	oldTable := m.table
	newTable := &Table{nodes: allocate(oldTable.len() * 2), ab: m.table.ab ^ 1}
	capacity := newTable.len()
	size := int64(0)
	for _, node := range oldTable.nodes {
		next := node.head
		for next != nil {
			next.next[newTable.ab], next.prev[newTable.ab] = nil, nil
			newNode := newTable.nodes[indexOf(next.hash, capacity)]
			if newNode.head == nil {
				newNode.head, newNode.tail = next, next
			} else {
				newNode.tail.next[newTable.ab], next.prev[newTable.ab], newNode.tail = next, newNode.tail, next
			}
			size++
			newNode.size++
			next = next.next[oldTable.ab]
		}
	}
	m.size = size
	m.table = newTable
}

func (m *HashMap) getNodeEntry(t *Table, n *Node, k interface{}) *Entry {
	next := n.head
	for next != nil {
		if next.k == k && next.flag == 0 {
			return next
		}
		next = next.next[t.ab]
	}
	return nil
}

func (m *HashMap) Get(k interface{}) (interface{}, bool) {
	t := m.table
	n, _ := t.getKeyNode(k)
	e := m.getNodeEntry(t, n, k)
	if e != nil {
		return e.Value(), true
	}
	return nil, false
}

func (m *HashMap) Del(k interface{}) bool {
	m.RLock()
	defer m.RUnlock()

	t := m.table
	n, _ := t.getKeyNode(k)
	n.Lock()
	defer n.Unlock()
	if e := m.getNodeEntry(t, n, k); e != nil {
		if e.prev[t.ab] == nil && e.next[t.ab] == nil {
			n.head, n.tail = nil, nil
		} else if e.prev[t.ab] == nil {
			n.head, n.head.prev[t.ab] = e.next[t.ab], nil
		} else if e.next[t.ab] == nil {
			n.tail, n.tail.next[t.ab] = e.prev[t.ab], nil
		} else {
			e.prev[t.ab].next[t.ab], e.next[t.ab].prev[t.ab] = e.next[t.ab], e.prev[t.ab]
		}
		oldAb := t.ab ^ 1
		if e.prev[oldAb] != nil {
			e.prev[oldAb].next[oldAb] = e.next[oldAb]
		}
		if e.next[oldAb] != nil {
			e.next[oldAb].prev[oldAb] = e.prev[oldAb]
		}
		atomic.AddInt64(&n.size, -1)
		atomic.AddInt64(&m.size, -1)
		return true
	}
	return false
}

func (m *HashMap) LogicDel(k interface{}) bool {
	h, t := hash(k), m.table
	n := t.nodes[indexOf(h, t.len())]

	//If key exists
	if e := m.getNodeEntry(t, n, k); e != nil {
		if atomic.CompareAndSwapInt32(&e.flag, 0, 1) {
			atomic.AddInt64(&n.size, -1)
			atomic.AddInt64(&m.size, -1)
			return true
		}
	}
	return false
}

func (e *Entry) Value() interface{} {
	return *(*interface{})(e.p)
}

func (e *Entry) Key() interface{} {
	return e.k
}

func (e *Entry) Flag() int32 {
	return e.flag
}

func (m *HashMap) Foreach(fn func(e *Entry)) {
	t := m.table
	for _, node := range t.nodes {
		next := node.head
		for next != nil {
			fn(next)
			next = next.next[t.ab]
		}
	}
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
	t := m.table
	data := map[string]interface{}{}
	for _, node := range t.nodes {
		next := node.head
		for next != nil {
			data[fmt.Sprintf("%v", next.k)] = next.Value()
			next = next.next[t.ab]
		}
	}
	return json.Marshal(data)
}
