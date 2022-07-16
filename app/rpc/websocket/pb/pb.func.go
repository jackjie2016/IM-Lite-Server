package pb

var (
	defaultListUserConversationReqPageSize int64 = 5000
	defaultPullMsgBySeqReqPageSize         int64 = 200
)

func (x *ListUserConversationReq) GetPageSizeX() *int64 {
	if x.PageSize == 0 {
		return &defaultListUserConversationReqPageSize
	}
	ps := int64(x.PageSize)
	return &ps
}

// IsStore 是否存储
func (x *MsgData) IsStore() bool {
	if x.MsgOptions == nil {
		return false
	}
	return x.MsgOptions.Storage
}

// IsIncrUnread 是否增加未读数
func (x *MsgData) IsIncrUnread() bool {
	if x.MsgOptions == nil {
		return false
	}
	return x.MsgOptions.Unread
}

func (x *PullMsgBySeqReq) GetPageSizeX() int64 {
	if x.PageSize == 0 {
		return defaultPullMsgBySeqReqPageSize
	}
	ps := int64(x.PageSize)
	return ps
}
