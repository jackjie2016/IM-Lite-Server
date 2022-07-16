// Code generated by goctl. DO NOT EDIT.
package types

type ReqGetConversationList struct {
}

type ModelConversation struct {
	ConvId      string `json:"convId"`
	MaxSeq      uint32 `json:"maxSeq"`
	MinSeq      uint32 `json:"minSeq"`
	UnreadCount uint32 `json:"unreadCount"`
}

type RespGetConversationList struct {
	Conversations []*ModelConversation `json:"conversations"`
}

type ReqPullMessage struct {
	ConvId    string   `json:"convId"`
	SeqList   []uint32 `json:"seqList"`
	OldestSeq uint32   `json:"oldestSeq,optional"`
	PageSize  int32    `json:"pageSize,optional"`
}

type ReqPullMessages struct {
	Convs []*ReqPullMessage `json:"convs"`
}

type RespPullMessages struct {
	MessageLists []byte `json:"messages"`
}

type ReqSendMessage struct {
	Message []byte `json:"message"`
}

type RespSendMessage struct {
	FailedMsg string `json:"failedMsg"`
}

type ReqClearUnreadCount struct {
	ConvId string `json:"convId"`
}

type RespClearUnreadCount struct {
}

type ReqClearHistory struct {
	ConvId string `json:"convId"`
}

type RespClearHistory struct {
}