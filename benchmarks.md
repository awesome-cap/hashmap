# Benchmarks

```text
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkReadHashMapUint-12                       532046              2456 ns/op
BenchmarkReadGoMapUintUnsafe-12                   518766              2227 ns/op
BenchmarkReadGoMapUintMutex-12                     33429             35621 ns/op
BenchmarkReadGoSyncMapUint-12                     212886              5884 ns/op
BenchmarkReadHashMapWithWritesUint-12             267873              4338 ns/op
BenchmarkReadGoMapWithWritesUintMutex-12            7822            132522 ns/op
BenchmarkReadGoSyncMapWithWritesUint-12           167682              6816 ns/op
BenchmarkWriteHashMapUint-12                       12524             84872 ns/op
BenchmarkWriteGoMapUnsafeUint-12                   63816             19729 ns/op
BenchmarkWriteGoMapMutexUint-12                    30558             34013 ns/op
BenchmarkWriteGoSyncMapUint-12                      9079            122624 ns/op
BenchmarkStrconv-12                             35964432                28.46 ns/op
BenchmarkSingleInsertAbsent-12                   2504090               601.5 ns/op
BenchmarkSingleInsertAbsentSyncMap-12            1355524               752.4 ns/op
BenchmarkSingleInsertPresent-12                 16673116                62.34 ns/op
BenchmarkSingleInsertPresentSyncMap-12          12789591                88.88 ns/op
BenchmarkMultiInsertDifferent-12                  520726              2517 ns/op
BenchmarkMultiInsertDifferentSyncMap-12           250897              4622 ns/op
BenchmarkMultiInsertSame-12                      1000000              1017 ns/op
BenchmarkMultiInsertSameSyncMap-12                335350              3646 ns/op
BenchmarkMultiGetSame-12                         4290022               302.4 ns/op
BenchmarkMultiGetSameSyncMap-12                  3872754               320.0 ns/op
BenchmarkMultiGetSetDifferent-12                  370822              2722 ns/op
BenchmarkMultiGetSetDifferentSyncMap-12           190112              7079 ns/op
BenchmarkMultiGetSetBlockSyncMap-12              1811307               673.4 ns/op
```

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