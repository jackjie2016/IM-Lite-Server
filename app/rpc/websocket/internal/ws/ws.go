package ws

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type UserConn struct {
	*websocket.Conn
}

func NewUserConn(conn *websocket.Conn) *UserConn {
	return &UserConn{Conn: conn}
}

var ws *Ws

type Ws struct {
	rwLock         sync.RWMutex
	svcCtx         *svc.ServiceContext
	WsConnToUser   map[*UserConn]map[string]string
	WsUserToConn   map[string]map[string]*UserConn
	allPlatformMap sync.Map
}

func GetWs(svcCtx *svc.ServiceContext) *Ws {
	if ws == nil {
		ws = &Ws{
			svcCtx: svcCtx,
		}
		ws.WsUserToConn = make(map[string]map[string]*UserConn)
		ws.WsConnToUser = make(map[*UserConn]map[string]string)
		go ws.checkUserConn(time.Second * 60)
	}
	return ws
}

func (l *Ws) WsUpgrade(ctx context.Context, uid string, req *ConnRequest, w http.ResponseWriter, r *http.Request) error {
	logger := logx.WithContext(ctx)
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode:      websocket.CompressionDisabled,
		CompressionThreshold: 0,
		OriginPatterns:       []string{"*"},
	})
	if err != nil {
		logger.Errorf("WsUpgrade error: %s", err)
		return err
	}
	newConn := NewUserConn(conn)
	err = l.AddUserConn(uid, req.Platform, newConn)
	if err != nil {
		logger.Errorf("AddUserConn error: %s", err)
		return err
	}
	return nil
}
