# expgo/sync

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/expgo/sync/blob/master/LICENSE)

`expgo/sync` 是一个Go语言项目，提供高级同步原语和并发控制机制。这个项目旨在扩展Go标准库中的`sync`包，提供更多功能和更灵活的并发编程工具。

## 安装

要安装`expgo/sync`，请确保你的Go环境版本在1.11以上，并运行以下命令：

```bash
go get github.com/expgo/sync
```

## 功能

`sync.Once`接口

* 扩展了 Go 标准库中的 `sync.Once`
* 原子变量 `done` 来标记函数已经被执行过，如果已经被执行过直接返回 `nil`。

接口提供三种调用方式：
```go
type Once interface {
	// 与系统提供的sync的Once提供相同的功能，区别在于系统的Once的f函数不提供返回错误
	Do(f func() error) error
	// 提供执行超时的功能
	DoTimeout(timeout time.Duration, f func() error) error
	// 提供支持context.Context的功能
	DoContext(ctx context.Context, f func() error) error
}
```

`sync.Mutex`接口提供和系统包相同的操作方式。

`NewMutex`支持在`系统Mutex`，`loggedMutex`和`deadlock`进行切换，默认为系统`Mutex`。

* 通过环境变量`SYNC_USE_SLOW_LOCK`，切换到`loggedMutex`，支持将超过100ms的锁打印输出
* 通过环境变量`SYNC_USE_DEADLOCK`，切换到`github.com/sasha-s/go-deadlock`

`sync.RWMutex`接口提供和系统包相同的操作方式。

`NewRWMutex`支持在`系统RWMutex`，`loggedRWMutex`和`deadlock`进行切换，默认为系统`RWMutex`。

* 通过环境变量`SYNC_USE_SLOW_LOCK`，切换到`loggedRWMutex`，支持将超过100ms的锁打印输出
* 通过环境变量`SYNC_USE_DEADLOCK`，切换到`github.com/sasha-s/go-deadlock`

## 许可证

`expgo/sync` 是根据MIT许可证发布的。有关详细信息，请参阅[LICENSE](https://github.com/expgo/sync/blob/master/LICENSE)文件。