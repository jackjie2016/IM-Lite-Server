package config

import (
	"github.com/IM-Lite/IM-Lite-Server/common/xredis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Redis        xredis.Config
	WebsocketRpc zrpc.RpcClientConf
}
