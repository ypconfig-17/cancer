# Cancer

用 Go 实现的 STOMP 消息服务实验。项目将网络连接、协议帧处理、发布订阅和本地队列拆开，便于沿消息路径阅读。

| 入口 | 用途 |
| --- | --- |
| [cmd/main.go](cmd/main.go) | 创建并启动服务，示例监听端口为 9999 |
| [stomp](stomp/) | STOMP 帧的解析与表示 |
| [发布与订阅](pub.go) | 查看消息发布和订阅分发的实现 |
| [store](store/) | 本地文件队列及索引 |

## 一条消息如何流转

1. [服务入口](server.go) 注册 HTTP 处理，接入连接并把消息交给协议处理链。
2. [帧处理器](handler.go) 识别 STOMP 的 CONNECT、SUBSCRIBE、SEND 和 UNSUBSCRIBE 指令。
3. [发布](pub.go) 和[订阅](sub.go) 分别处理消息与订阅的分发，文件队列和索引由 [store](store/) 负责。

阅读时可以从 SEND 分支追到发布处理，再从 SUBSCRIBE 分支追到订阅注册。协议帧的表示与编解码集中在 [stomp](stomp/)，便于单独检查。

## 源码中的设计细节

- **连接与协议处理分离。** [server.go](server.go) 在 HTTP 入口接受 WebSocket 连接，为每条连接创建 Channel，并安装帧处理链。 [Channel](channel.go) · [Pipeline](pipe.go)
- **按目的地写队列。** 发布端按消息的目的地打开文件队列并追加消息。订阅侧通过文件队列读取任务，再由循环处理。 [发布实现](pub.go) · [队列读取](poller.go) · [订阅任务循环](looper.go)
- **本地运行边界。** 示例使用 WebSocket 路径 `/`、端口 `9999`，队列数据写入进程工作目录下的 `db`。服务端当前的握手回调直接放行，适合在隔离的本地环境阅读和试跑。 [启动入口](cmd/main.go) · [握手与存储配置](server.go)

## 验证

项目声明 Go 1.19。可运行 `go test ./stomp` 检查协议帧，也可用 `go run ./cmd` 启动示例服务。队列实验涉及工作目录中的文件，建议在独立目录验证。

当前 `go test ./stomp` 在编译测试时会报错。帧测试仍使用旧的构造参数和字符串消息体，实现已采用新的参数和字节消息体。 [帧测试](stomp/frame_test.go)
