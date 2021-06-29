package hashmap

import (
	"sync"
	"sync/atomic"
	"testing"
)

const benchmarkItemCount = 1 << 10 // 1024

func setupHashMap(b *testing.B) *HashMap {
	m := New()
	for i := uintptr(0); i < benchmarkItemCount; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()
	return m
}

func setupGoMap(b *testing.B) map[uintptr]uintptr {
	m := make(map[uintptr]uintptr)
	for i := uintptr(0); i < benchmarkItemCount; i++ {
		m[i] = i
	}

	b.ResetTimer()
	return m
}

func setupGoSyncMap(b *testing.B) *sync.Map {
	m := &sync.Map{}
	for i :=uintptr(0); i < benchmarkItemCount; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()
	return m
}

func BenchmarkReadHashMapUint(b *testing.B) {
	m := setupHashMap(b)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := uintptr(0); i < benchmarkItemCount; i++ {
				j, _ := m.Get(i)
				if j != i {
					b.Fail()
				}
			}
		}
	})
}

func BenchmarkReadGoMapUintUnsafe(b *testing.B) {
	m := setupGoMap(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i :=uintptr(0); i < benchmarkItemCount; i++ {
				j := m[i]
				if j != i {
					b.Fail()
				}
			}
		}
	})
}

func BenchmarkReadGoMapUintMutex(b *testing.B) {
	m := setupGoMap(b)
	l := &sync.RWMutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i :=uintptr(0); i < benchmarkItemCount; i++ {
				l.RLock()
				j := m[i]
				l.RUnlock()
				if j != i {
					b.Fail()
				}
			}
		}
	})
}

func BenchmarkReadGoSyncMapUint(b *testing.B) {
	m := setupGoSyncMap(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i :=uintptr(0); i < benchmarkItemCount; i++ {
				j, _ := m.Load(i)
				if j != i {
					b.Fail()
				}
			}
		}
	})
}

func BenchmarkReadHashMapWithWritesUint(b *testing.B) {
	m := setupHashMap(b)
	var writer int32

	b.RunParallel(func(pb *testing.PB) {
		// use 1 thread as writer
		if atomic.CompareAndSwapInt32(&writer, 0, 1) {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					m.Set(i, i)
				}
			}
		} else {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					j, _ := m.Get(i)
					if j != i {
						b.Fail()
					}
				}
			}
		}
	})
}

func BenchmarkReadGoMapWithWritesUintMutex(b *testing.B) {
	m := setupGoMap(b)
	l := &sync.RWMutex{}
	var writer int32

	b.RunParallel(func(pb *testing.PB) {
		// use 1 thread as writer
		if atomic.CompareAndSwapInt32(&writer, 0, 1) {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					l.Lock()
					m[i] = i
					l.Unlock()
				}
			}
		} else {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					l.RLock()
					j := m[i]
					l.RUnlock()
					if j != i {
						b.Fail()
					}
				}
			}
		}
	})
}


func BenchmarkReadGoSyncMapWithWritesUint(b *testing.B) {
	m := setupGoSyncMap(b)
	var writer int32

	b.RunParallel(func(pb *testing.PB) {
		// use 1 thread as writer
		if atomic.CompareAndSwapInt32(&writer, 0, 1) {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					m.Store(i, i)
				}
			}
		} else {
			for pb.Next() {
				for i :=uintptr(0); i < benchmarkItemCount; i++ {
					j, _ := m.Load(i)
					if j != i {
						b.Fail()
					}
				}
			}
		}
	})
}

func BenchmarkWriteHashMapUint(b *testing.B) {
	m := New()

	for n := 0; n < b.N; n++ {
		for i :=uintptr(0); i < benchmarkItemCount; i++ {
			m.Set(i, i)
		}
	}
}

func BenchmarkWriteGoMapUnsafeUint(b *testing.B) {
	m := make(map[uintptr]uintptr)

	for n := 0; n < b.N; n++ {
		for i :=uintptr(0); i < benchmarkItemCount; i++ {
			m[i] = i
		}
	}
}


func BenchmarkWriteGoMapMutexUint(b *testing.B) {
	m := make(map[uintptr]uintptr)
	l := &sync.RWMutex{}

	for n := 0; n < b.N; n++ {
		for i :=uintptr(0); i < benchmarkItemCount; i++ {
			l.Lock()
			m[i] = i
			l.Unlock()
		}
	}
}

func BenchmarkWriteGoSyncMapUint(b *testing.B) {
	m := &sync.Map{}

	for n := 0; n < b.N; n++ {
		for i :=uintptr(0); i < benchmarkItemCount; i++ {
			m.Store(i, i)
		}
	}
}
