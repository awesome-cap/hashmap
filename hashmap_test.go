package hashmap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewHashMap(t *testing.T) {
	hm := New()
	batch := 1000000
	start := time.Now().UnixNano()
	for i := 0; i < batch; i++ {
		hm.Set(i, i)
	}
	end := time.Now().UnixNano()
	fmt.Println("set:", (end-start)/1e6)

	start = time.Now().UnixNano()
	for i := 0; i < batch; i++ {
		v, e := hm.Get(i)
		if !e || v != i {
			t.Fatal("data err ", i)
		}
	}
	end = time.Now().UnixNano()
	fmt.Println("get:", (end-start)/1e6)
}

func TestHashMap_MSet(t *testing.T) {
	hm := New()
	batch := 1000000
	ks := make([]interface{}, batch)
	vs := make([]interface{}, batch)
	for i := 0; i < batch; i++ {
		ks[i] = i
		vs[i] = i
	}
	hm.MSet(ks, vs)
	for i := 0; i < batch; i++ {
		v, e := hm.Get(i)
		if !e || v != i {
			t.Fatal("data err ", i)
		}
	}
}

func TestHashMap_SetNil(t *testing.T) {
	m := New()
	m.Set("a", nil)
	_, ok := m.Get("a")
	if !ok {
		t.Fatal("set nil err")
	}
}

func TestHashMapCorrectness(t *testing.T) {
	hm := New()
	batch := 1000000
	start := time.Now().UnixNano()
	for i := 0; i < batch; i++ {
		hm.Set(i, i)
	}
	for i := 0; i < batch; i++ {
		v, _ := hm.Get(i)
		v2 := hm.Set(i, v.(int)+1)
		if v2 != i {
			t.Fatalf("v2 val : %v", v2)
		}
	}
	for i := 0; i < batch; i++ {
		v, e := hm.Get(i)
		if !e || v != i+1 {
			t.Fatal("data err ", i)
		}
	}
	end := time.Now().UnixNano()
	fmt.Println("get:", (end-start)/1e6)

	for i := 0; i < batch/2; i++ {
		b := hm.Del(i)
		if !b {
			t.Fatal("del err", i)
		}
	}

	for i := batch / 2; i < batch; i++ {
		b := hm.LogicDel(i)
		if !b {
			t.Fatal("del err", i)
		}
	}
	for i := batch / 2; i < batch; i++ {
		b := hm.LogicDel(i)
		if b {
			t.Fatal("del err", i)
		}
	}
	if hm.size != 0 {
		t.Fatal("size err")
	}
	for i := 0; i < batch; i++ {
		_, b := hm.Get(i)
		if b {
			t.Fatal("Get deleted value err", i)
		}
	}
}

func TestNewHashMap_Sync(t *testing.T) {
	hm := New()
	batch := 100000
	wg := sync.WaitGroup{}
	wg.Add(batch)
	for i := 0; i < batch; i++ {
		n := i
		go func() {
			hm.Set(strconv.Itoa(n), n)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(hm.size)
	if hm.size != int64(batch) {
		t.Fatal("TestNewHashMap_Sync SET ERR")
	}

	wg.Add(batch / 2)
	for i := 0; i < batch/2; i++ {
		n := i
		go func() {
			hm.Del(strconv.Itoa(n))
			wg.Done()
		}()
	}
	wg.Wait()
	if hm.size != int64(batch-batch/2) {
		t.Fatal("TestNewHashMap_Sync DEL ERR")
	}
}

func TestHashMap_MarshalJSON(t *testing.T) {
	m := New()
	m.Set("abc", "haha")
	m.Set(1, 2)
	m.Set("m", map[string]string{
		"hello": "world",
	})
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", b)
}

func TestHashMap_UnmarshalJSON(t *testing.T) {
	jsonStr := "{\"1\":2,\"abc\":\"haha\",\"m\":{\"hello\":\"world\"}}"
	m := New()
	err := json.Unmarshal([]byte(jsonStr), m)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	fmt.Println(m.Get("1"))
}

func TestHashMap_SetNX(t *testing.T) {
	m := New()
	b := m.SetNX("a", "b")
	assertEqual(t, b, true)
	b = m.SetNX("a", "b")
	assertEqual(t, b, false)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatal(fmt.Sprintf("%v not equal %v", a, b))
	}
}
