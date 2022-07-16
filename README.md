# IM-Lite-Server

轻量级 分布式的IM服务器

## 通讯流程

![通讯流程.svg](https://raw.githubusercontent.com/showurl/images/73a15955a98f2b6c2a21ab5be58d7ee1a89fdacc/IM-Lite%E9%80%9A%E8%AE%AF%E6%B5%81%E7%A8%8B.svg)

## 配置

### websocket.yaml

> 以下是可修改的配置

| 参数                      | 参数名称                      | 参数类型     | 参数描述                                    |
|-------------------------|---------------------------|----------|-----------------------------------------|
| ListenOn                | 监听地址                      | string   | rpc的监听地址                                |
| Mode                    | 模式                        | string   | 服务运行模式,可选 `dev`,`test`,`rt`,`pre`,`pro` |
| Websocket.Port          | websocket端口号              | int      | websocket端口号                            |
| Websocket.MaxConnNum    | websocket最大连接数            | int      | websocket最大连接数                          |
| Log.Mode                | 日志模式                      | string   | 可选 `console`, `file`, `volume`          |
| Log.Encoding            | 日志编码格式                    | string   | 可选 `json`, `plain`                      |
| Log.TimeFormat          | 日志时间格式                    | string   |                                         |
| Log.Path                | 日志保存路径                    | string   | 默认 `logs`                               |
| Log.Level               | 日志级别                      | string   | 可选 `info`, `error`, `severe`            |
| Log.Compress            | 日志压缩                      | bool     |                                         |
| Telemetry.Name          | 链路追踪名称                    | string   |                                         |
| Telemetry.Endpoint      | 链路追踪地址                    | string   |                                         |
| Telemetry.Sampler       | 链路追踪取样率                   | float64  | 链路追踪取样率 默认全取 `1.0`                      |
| Telemetry.Batcher       | 链路追踪器                     | string   | 支持 `jaeger`, `zipkin`                   |
| Redis.Host              | redis地址                   | string   | `localhost:6379`                        |
| Redis.Pass              | redis密码                   | string   |                                         |
| Redis.DB                | redis数据库                  | int      |                                         |
| Mongo.Uri               | mongo uri                 | string   | `mongodb://localhost/admin`             |
| Mongo.Database          | mongo数据库                  | string   | `im_lite`                               |
| Mongo.DBTimeout         | 请求超时时间，单位秒                | int      |                                         |
| Pushers.Msg.Addrs       | 用户发送消息存入kafka的配置，broker地址 | []string |                                         |
| Pushers.Msg.Topic       | 用户发送消息存入kafka的配置，topic名称  | string   |                                         |
| Pushers.Msg.User      | kafka sasl plain username | string   | 可不传                                     |
| Pushers.Msg.Passwd    | kafka sasl plain password | string   | 可不传                                     |
| Pushers.Msg.Partition    | topic 分区数                 | int      | 默认 `1`                                    |
| Pushers.Push.Addrs      | 消息落库后推送到kafka的配置，broker地址 | []string |                                         |
| Pushers.Push.Topic      | 消息落库后推送到kafka的配置，topic名称  | string   |                                         |
| Pushers.Push.User      | kafka sasl plain username | string   | 可不传                                     |
| Pushers.Push.Passwd    | kafka sasl plain password | string   | 可不传                                     |
| Pushers.Push.Partition    | topic 分区数 | string   | 默认 `1`                                     |
| Consumers.Msg.Addrs     | 消费用户发送到消息，broker地址        | []string |                                         |
| Consumers.Msg.Topic     | 消费用户发送到消息，topic名称         | string   |                                         |
| Consumers.Msg.Consumers | 启动的消费者数量                  | int      | 默认 `8`                                  |
| Consumers.Msg.User      | kafka sasl plain username | string   | 可不传                                     |
| Consumers.Msg.Passwd    | kafka sasl plain password | string   | 可不传                                     |
| Consumers.Push.Addrs    | 消费推送任务(广播模式)，broker地址     | []string |                                         |
| Consumers.Push.Topic    | 消费推送任务(广播模式)，topic名称      | string   |                                         |
| Consumers.Push.Consumers | 启动的消费者数量                  | int      | 默认 `8`                                  |
| Consumers.Push.User      | kafka sasl plain username | string   | 可不传                                     |
| Consumers.Push.Passwd    | kafka sasl plain password | string   | 可不传                                     |

### 环境变量

> IM-Lite-Server是分布式的，当启动多个Pod时，需要传递以下环境变量

- POD_NAME - Pod的名称 `POD_NAME会当作 push consumer 的 group id`
- POD_IP - Pod的IP

## docker-compose

## 编译安装

### websocket-rpc
> 由于IM-Lite-Server使用了 [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go) ，因此编译 `websocket-rpc`时 需要启用CGO
> 
> 建议直接在 docker 容器中编译，不要在本地编译。请参考 [Dockerfile](https://github.com/IM-Lite/IM-Lite-Server/blob/master/deploy/build/websocket-rpc/Dockerfile) 文件

