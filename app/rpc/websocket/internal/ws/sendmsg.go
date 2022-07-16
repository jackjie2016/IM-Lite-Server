package ws

import (
	"context"
	"nhooyr.io/websocket"
)

func (l *Ws) SendMsgToUidIgnoreErr(ctx context.Context, uid string, value []byte) {
	for _, conn := range l.GetUserConns(uid) {
		go l.sendMsg(ctx, conn, value)
	}
}

func (l *Ws) sendMsg(ctx context.Context, conn *UserConn, msg []byte) error {
	return conn.Write(ctx, websocket.MessageBinary, msg)
}
