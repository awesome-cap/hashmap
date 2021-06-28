# Introduce
Atomic HashMap.
# Install
```powershell
go get -u github.com/awesome-cap/hashmap
```
# Usage
```golang
import "github.com/awesome-cap/hashmap"

m := hashmap.New()
m.Set("hello", "world")
v, ok := m.Get("hello") // v: "world", ok: true
suc := m.Del("hello")   // suc: true
```
# Benchmark Test
```text
BenchmarkReadHashMapUint-4                        222816              5936 ns/op
BenchmarkReadGoMapUintUnsafe-4                    231386              5774 ns/op
BenchmarkReadGoMapUintMutex-4                      37956             31414 ns/op
BenchmarkReadGoSyncMapUint-4                       78129             14272 ns/op
BenchmarkReadHashMapWithWritesUint-4               89791             13368 ns/op
BenchmarkReadGoMapWithWritesUintMutex-4             6482            198641 ns/op
BenchmarkReadGoSyncMapWithWritesUint-4             71194             16693 ns/op
BenchmarkWriteHashMapUint-4                        23698             50918 ns/op
BenchmarkWriteGoMapUnsafeUint-4                    62028             19093 ns/op
BenchmarkWriteGoMapMutexUint-4                     34281             35057 ns/op
BenchmarkWriteGoSyncMapUint-4                       8593            132253 ns/op
BenchmarkStrconv-4                              38813596                31.6 ns/op
BenchmarkSingleInsertAbsent-4                    2265174               670 ns/op
BenchmarkSingleInsertAbsentSyncMap-4             1736335               825 ns/op
BenchmarkSingleInsertPresent-4                  35389147                34.3 ns/op
BenchmarkSingleInsertPresentSyncMap-4           12084397               101 ns/op
BenchmarkMultiInsertDifferent-4                   286462              3579 ns/op
BenchmarkMultiInsertDifferentSyncMap-4            308533              7283 ns/op
BenchmarkMultiInsertSame-4                       1000000              2180 ns/op
BenchmarkMultiInsertSameSyncMap-4                 414916              3367 ns/op
BenchmarkMultiGetSame-4                          3356172               393 ns/op
BenchmarkMultiGetSameSyncMap-4                   3089238               409 ns/op
BenchmarkMultiGetSetDifferent-4                   601536              2701 ns/op
BenchmarkMultiGetSetDifferentSyncMap-4            207987              7206 ns/op
BenchmarkMultiGetSetBlockSyncMap-4               1000000              1299 ns/op
```