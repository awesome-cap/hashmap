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
# Benchmarks

[Benchmark](./benchmarks.md)