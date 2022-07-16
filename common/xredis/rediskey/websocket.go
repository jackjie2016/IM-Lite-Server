package rediskey

import "fmt"

func UserConn(uid, platformId string, podip string) string {
	return fmt.Sprintf("websocket:userconn:%s:%s", uid, platformId)
}

func ConversationSeq(conversationID string) string {
	return fmt.Sprintf("websocket:conversationseq:%s", conversationID)
}

func ConversationUnread(belongUser string) string {
	return fmt.Sprintf("websocket:conversationunread:%s", belongUser)
}

func UserSubscribedConversations(conversationId, uid string) string {
	return fmt.Sprintf("websocket:usersubscribedconversations:%s:%s", conversationId, uid)
}
