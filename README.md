# DetectLock-go

用于检测运行时死锁。获得产生死锁的goroutine编号后，结合go-pprof的堆栈信息，即可找到代码相关位置。

goroutine的锁信息记录，采用分片锁，降低了多线程并发竞争的性能影响。

- 支持 `sync.Mutex`、`sync.RWMutex`；

- 支持检测 `锁重入`、`多把锁顺序不一致导致的互锁`；

- 支持 启用 / 禁用 调试功能。

## 用法

```go
import (
    "github.com/berkaroad/detectlock-go/detectlock"
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

`goroutine <协程ID>: [(<锁标识>, <锁状态>), ...]`

锁标识：相同标识即为同一把锁。

锁状态，共4种：

- acquired

  获得Mutex锁，或RWMutex写锁。

- wait

  等待Mutex锁，或等待RWMutex写锁。

- r-acquired

  获得RWMutex读锁。

- r-wait

  等待RWMutex读锁。

## 检测到死锁的示例

- 检测 sync.Mutex 多把锁顺序不一致导致的互锁

```plain
--- DetectAcquired ---
goroutine 53: [(0xc000014080, acquired), (0xc000014088, wait)]
goroutine 54: [(0xc000014088, acquired), (0xc000014080, wait)]

--- DetectLockedEachOther ---
goroutine 53: [(0xc000014080, acquired), (0xc000014088, wait)]
goroutine 54: [(0xc000014088, acquired), (0xc000014080, wait)]

```

- 检测 sync.Mutex 锁重入

```plain
--- DetectAcquired ---
goroutine 15: [(0xc0000140a0, acquired), (0xc0000140a0, wait)]

--- DetectReentry ---
goroutine 15: [(0xc0000140a0, acquired), (0xc0000140a0, wait)]

```

- 检测 sync.Mutex、sync.RWMutex 多把锁顺序不一致导致的互锁

```plain
--- DetectAcquired ---
goroutine 9: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 11: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 13: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 15: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 49: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 51: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 53: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 56: [(0xc0000ae028, acquired), (0xc0000180c0, wait)]
goroutine 57: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 59: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]

--- DetectLockedEachOther ---
goroutine 9: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 11: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 13: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 15: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 49: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 51: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 53: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 56: [(0xc0000ae028, acquired), (0xc0000180c0, wait)]
goroutine 57: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]
goroutine 59: [(0xc0000180c0, r-acquired), (0xc0000ae028, wait)]

```

## 发布版本

### v1.0 (2021-10-31)

- 支持 `sync.Mutex`、`sync.RWMutex`；

- 支持检测 `锁重入`、`多把锁顺序不一致导致的互锁`；

- 支持 启用 / 禁用 调试功能。
