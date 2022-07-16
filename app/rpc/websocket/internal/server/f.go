package server

import (
	"context"
	"fmt"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/logic"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/ws"
	"log"
	"net/http"
)

var serverCtx = context.Background()

func (s *WebsocketServiceServer) Websocket() {
	http.HandleFunc("/", ws.Handler(s.svcCtx))
	wsAddr := fmt.Sprintf(":%d", s.svcCtx.Config.Websocket.Port)
	log.Printf("ws listen on %s...", wsAddr)
	log.Fatal(http.ListenAndServe(wsAddr, nil))
}

func (s *WebsocketServiceServer) Consumers() {
	logic.NewConsumeLogic(serverCtx, s.svcCtx)
}
