package _struct

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewHashMap(t *testing.T) {
	hm := New()
	batch := 1000000
	start := time.Now().UnixNano()
	for i := 0; i < batch; i ++ {
		hm.Set(i, i)
	}
	end := time.Now().UnixNano()
	fmt.Println("set:", (end - start) / 1e6)

	start = time.Now().UnixNano()
	for i := 0; i < batch; i ++ {
		v, e := hm.Get(i)
		if ! e || v != i{
			log.Fatal("data err ", i)
		}
	}
	end = time.Now().UnixNano()
	fmt.Println("get:", (end - start) / 1e6)
}

func TestNewHashMap_Sync(t *testing.T) {
	hm := New()
	batch := 100000
	wg := sync.WaitGroup{}
	wg.Add(batch)
	for i := 0; i < batch; i ++ {
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
	for i := 0; i < batch / 2; i ++ {
		n := i
		go func() {
			hm.Del(strconv.Itoa(n))
			wg.Done()
		}()
	}
	wg.Wait()
	if hm.size != int64(batch - batch / 2) {
		t.Fatal("TestNewHashMap_Sync DEL ERR")
	}
}
