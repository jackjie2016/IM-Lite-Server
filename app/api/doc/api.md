### 1. "获取会话列表"

1. route definition

- Url: /im/v1/conversation/list
- Method: POST
- Request: `ReqGetConversationList`
- Response: `RespGetConversationList`

2. request definition



```golang
type ReqGetConversationList struct {
}
```


3. response definition



```golang
type RespGetConversationList struct {
	Message []byte `json:"message"`
}
```

### 2. "批量拉取消息"

1. route definition

- Url: /im/v1/message/pull/batch
- Method: POST
- Request: `ReqPullMessages`
- Response: `RespPullMessages`

2. request definition



```golang
type ReqPullMessages struct {
	Message []byte `json:"message"`
}
```


3. response definition



```golang
type RespPullMessages struct {
	MessageLists []byte `json:"messages"`
}
```

### 3. "发送消息"

1. route definition

- Url: /im/v1/message/send
- Method: POST
- Request: `ReqSendMessage`
- Response: `RespSendMessage`

2. request definition



```golang
type ReqSendMessage struct {
	Message []byte `json:"message"`
}
```


3. response definition



```golang
type RespSendMessage struct {
	FailedMsg string `json:"failedMsg"`
}
```

### 4. "清空未读数量"

1. route definition

- Url: /im/v1/conversation/clear/unread
- Method: POST
- Request: `ReqClearUnreadCount`
- Response: `RespClearUnreadCount`

2. request definition



```golang
type ReqClearUnreadCount struct {
	ConvId string `json:"convId"`
}
```


3. response definition



```golang
type RespClearUnreadCount struct {
}
```

### 5. "清空会话消息"

1. route definition

- Url: /im/v1/conversation/clear/history
- Method: POST
- Request: `ReqClearHistory`
- Response: `RespClearHistory`

2. request definition



```golang
type ReqClearHistory struct {
	ConvId string `json:"convId"`
}
```


3. response definition



```golang
type RespClearHistory struct {
}
```

