package svc

import (
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/websocketservice"
	"github.com/zeromicro/go-zero/zrpc"
)

func (s *ServiceContext) WebsocketService() websocketservice.WebsocketService {
	if s.websocketService == nil {
		s.websocketService = websocketservice.NewWebsocketService(zrpc.MustNewClient(s.Config.WebsocketRpc))
	}
	return s.websocketService
}
