package database

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/model"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewDefault(svcCtx *svc.ServiceContext, ctx context.Context) Database {
	return &Default{svcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx)}
}

type Database interface {
	// HasConversation 是否有这个会话
	HasConversation(conversationID string) (bool, error)
	// InsertMayUpdateOneMessage 新增消息 如果 clientMsgID 重复则更新原有消息
	InsertMayUpdateOneMessage(chatLog *model.ChatLog) error
	// CreateConversation 新建会话
	CreateConversation(
		conversation *model.Conversation,
	) (string, error)
	// DeleteConversationAllMessages 删除会话所有消息
	DeleteConversationAllMessages(conversationID, userID string) error
	// DeleteConversation 删除会话
	DeleteConversation(conversationID string) error
	// ListUserConversations 获取用户的会话列表
	ListUserConversations(userID string, pageNo int, pageSize int64) ([]*model.MidUserConversation, error)
	// PullMsgBySeq 拉取会话消息
	PullMsgBySeq(conversationID string, userID string, maxSeq uint32, pageSize int64) ([]*model.ChatLog, error)
	IncrUnread(conversation *model.Conversation, msg *pb.MsgData)
	GetConversation(conversationID string) (*model.Conversation, error)
}
