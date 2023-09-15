# DetectLock-go

用于检测运行时死锁。获得产生死锁的goroutine编号后，结合go-pprof的堆栈信息，即可找到代码相关位置。

goroutine的锁信息记录，采用分片锁，降低了多线程并发竞争的性能影响。

- 支持 `sync.Mutex`、`sync.RWMutex`

- 支持检测 `锁重入`、`多把锁顺序不一致导致的互锁`

- 支持 启用 / 禁用 调试功能。

## 用法

```go
import (
    "github.com/berkaroad/detectlock-go"
)

// 应用启动时，设置启用调试
detectlock.EnableDebug()

// 声明 sync.Mutex、sync.RWMutex 替换为 detectlock.Mutex、detectlock.RWMutex
var locker1 *detectlock.Mutex = &detectlock.Mutex{}
var locker2 *detectlock.RWMutex = &detectlock.RWMutex{}

// 异步检测死锁
items := detectlock.Items()
fmt.Println(detectlock.DetectAcquired(items)) // 检测获得锁的goroutine列表
fmt.Println(detectlock.DetectReentry(items)) // 检测锁重入的goroutine列表
fmt.Println(detectlock.DetectLockedEachOther(items)) // 检测互锁的goroutine列表

// 关闭调试，并清理锁使用信息
detectlock.DisableDebug()
```

## 数据格式

`goroutine <协程ID>: [(<锁标识>, <锁状态>, <调用者函数>(file: <源码文件名>:<源码行号>)), ...]`

- 锁标识：相同标识即为同一把锁。

- 锁状态，共4种：

  - acquired

    获得Mutex锁，或RWMutex写锁。

  - wait

    等待Mutex锁，或等待RWMutex写锁。

  - r-acquired

    获得RWMutex读锁。

  - r-wait

    等待RWMutex读锁。

- 调用者函数

  调用了加锁操作的调用者的函数名，包含所属的包名。

- 源码文件名

  调用了加锁操作的调用者的函数所在的源文件名。

- 源码行号

  调用了加锁操作的调用者的函数所在的源文件中的行号。

## 检测到死锁的示例

- 检测 sync.Mutex 多把锁顺序不一致导致的互锁

```plain
--- DetectAcquired ---
goroutine 29: [(0xc0000b4008, acquired, main.B (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:30)), (0xc0000b4000, wait, main.B (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:33))]
goroutine 30: [(0xc0000b4000, acquired, main.A (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:20)), (0xc0000b4008, wait, main.A (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:23))]

--- DetectLockedEachOther ---
goroutine 29: [(0xc0000b4008, acquired, main.B (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:30)), (0xc0000b4000, wait, main.B (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:33))]
goroutine 30: [(0xc0000b4000, acquired, main.A (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:20)), (0xc0000b4008, wait, main.A (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex01/main.go:23))]

```

- 检测 sync.Mutex 锁重入

```plain
--- DetectAcquired ---
goroutine 53: [(0xc0000160b8, acquired, main.C (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex02/main.go:19)), (0xc0000160b8, wait, main.C (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex02/main.go:22))]

--- DetectReentry ---
goroutine 53: [(0xc0000160b8, acquired, main.C (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex02/main.go:19)), (0xc0000160b8, wait, main.C (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex02/main.go:22))]

```

- 检测 sync.Mutex、sync.RWMutex 多把锁顺序不一致导致的互锁

```plain
--- DetectAcquired ---
goroutine 8: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 10: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 13: [(0xc0000160b8, acquired, main.E (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:30)), (0xc0000180c0, wait, main.E (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:33))]
goroutine 14: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 16: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 50: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 52: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 54: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 56: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 58: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]

--- DetectLockedEachOther ---
goroutine 8: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 10: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 13: [(0xc0000160b8, acquired, main.E (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:30)), (0xc0000180c0, wait, main.E (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:33))]
goroutine 14: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 16: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 50: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 52: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 54: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 56: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]
goroutine 58: [(0xc0000180c0, r-acquired, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:21)), (0xc0000160b8, wait, main.D (file: /home/berkaroad/github/berkaroad/detectlock-go/examples/mutex03/main.go:24))]

```

## Benchmark

通过对 `启用调试` 、`禁用调试` 、 `原生` 的锁操作的性能对比，`启用调试` 模式下性能比较差，因此不适合用于生产环境。而对于 `禁用调试` 、 `原生` 的对比，性能损耗相差无几。

```sh
make go-bench
```

```plain
GO111MODULE=on GOPROXY=https://goproxy.cn,direct go test -run=none -count=1 -benchtime=10000x -benchmem -bench=. ./... | grep Benchmark
BenchmarkMutex_Lock/EnableDebug-4                  10000             97118 ns/op            5928 B/op         75 allocs/op
BenchmarkMutex_Lock/DisableDebug-4                 10000               462.6 ns/op           160 B/op         10 allocs/op
BenchmarkMutex_Lock/sync.Mutex-4                   10000               401.7 ns/op           160 B/op         10 allocs/op
BenchmarkMutex_TryLock/EnableDebug-4               10000            103766 ns/op            5928 B/op         75 allocs/op
BenchmarkMutex_TryLock/DisableDebug-4              10000               447.7 ns/op           160 B/op         10 allocs/op
BenchmarkMutex_TryLock/sync.Mutex-4                10000               391.4 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_RLock/EnableDebug-4               10000            102587 ns/op            5928 B/op         75 allocs/op
BenchmarkRWMutex_RLock/DisableDebug-4              10000               603.0 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_RLock/sync.RWMutex-4              10000               415.0 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_TryRLock/EnableDebug-4            10000            413464 ns/op            6090 B/op         70 allocs/op
BenchmarkRWMutex_TryRLock/DisableDebug-4           10000               407.6 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_TryRLock/sync.RWMutex-4           10000               428.3 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_Lock/EnableDebug-4                10000            101591 ns/op            5928 B/op         75 allocs/op
BenchmarkRWMutex_Lock/DisableDebug-4               10000               544.3 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_Lock/sync.RWMutex-4               10000               875.6 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_TryLock/EnableDebug-4             10000            102170 ns/op            5928 B/op         75 allocs/op
BenchmarkRWMutex_TryLock/DisableDebug-4            10000               471.0 ns/op           160 B/op         10 allocs/op
BenchmarkRWMutex_TryLock/sync.RWMutex-4            10000               456.1 ns/op           160 B/op         10 allocs/op

```

## 发布版本

### v1.1 (2023-09-15)

- 加锁时，额外收集调用者函数栈

  用于问题诊断时，可以看到调用者的函数完整名、代码行。

- 补全 `sync.Mutex`、`sync.RWMutex`

  补充缺失的函数： `TryLock()`、`TryRLock()`。

- 优化了未开启调试场景下，默认不占用10Mb内存

### v1.0 (2021-10-31)

- 支持 `sync.Mutex`、`sync.RWMutex`；

- 支持检测 `锁重入`、`多把锁顺序不一致导致的互锁`；

- 支持 启用 / 禁用 调试功能。
