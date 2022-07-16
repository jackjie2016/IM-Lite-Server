package ws

import (
	"context"
	"fmt"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"nhooyr.io/websocket"
	"time"
)

func (l *Ws) AddUserConn(uid string, platformID string, conn *UserConn) error {
	if l.svcCtx.Config.Websocket.MaxConnNum > 0 && len(l.WsConnToUser) >= l.svcCtx.Config.Websocket.MaxConnNum {
		// 连接数超过限制
		return ErrMaxConnNum
	}
	// 加入redis
	err := l.svcCtx.Redis().HSet(
		context.Background(), rediskey.UserConn(uid, platformID, l.svcCtx.Config.PodIp()),
		fmt.Sprintf("%s%s", l.svcCtx.Config.PodIp(), l.svcCtx.Config.ListenOn), // key
		time.Now().UnixMilli(),                                                 // value
	).Err()
	if err != nil {
		logx.Error("add user conn to redis err", " uid ", uid, " platform ", platformID, " err ", err)
		return err
	} else {
		logx.Info("add user conn to redis success", " uid ", uid, " platform ", platformID)
	}
	if oldConnMap, ok := l.WsUserToConn[uid]; ok {
		oldConnMap[platformID] = conn
		l.WsUserToConn[uid] = oldConnMap
	} else {
		i := make(map[string]*UserConn)
		i[platformID] = conn
		l.WsUserToConn[uid] = i
	}
	if oldStringMap, ok := l.WsConnToUser[conn]; ok {
		oldStringMap[platformID] = uid
		l.WsConnToUser[conn] = oldStringMap
	} else {
		i := make(map[string]string)
		i[platformID] = uid
		l.WsConnToUser[conn] = i
	}
	l.allPlatformMap.Store(platformID, struct{}{})
	return nil
}

func (l *Ws) GetUserConn(uid string, platform string) *UserConn {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	if connMap, ok := l.WsUserToConn[uid]; ok {
		if conn, flag := connMap[platform]; flag {
			return conn
		}
	}
	return nil
}

func (l *Ws) GetUserConns(uid string) (conns []*UserConn) {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	if connMap, ok := l.WsUserToConn[uid]; ok {
		for _, conn := range connMap {
			conns = append(conns, conn)
		}
	}
	return conns
}

func (l *Ws) DelUserConn(ctx context.Context, conn *UserConn) error {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	uid, platform := l.GetUserUid(conn)
	err := l.svcCtx.Redis().HDel(
		ctx, rediskey.UserConn(uid, platform, l.svcCtx.Config.PodIp()),
		fmt.Sprintf("%s%s", l.svcCtx.Config.PodIp(), l.svcCtx.Config.ListenOn), // key
	).Err()
	if err != nil {
		logx.Error("del user conn from redis err", " uid ", uid, " platform ", platform, " err ", err)
		return err
	}

	// 是否需要删除用户的所有连接
	go func() {
		result, err := l.svcCtx.Redis().HLen(context.Background(), rediskey.UserConn(uid, platform, l.svcCtx.Config.PodIp())).Result()
		if err != nil {
			logx.Error("get user conn num from redis err", " uid ", uid, " platform ", platform, " err ", err)
		} else {
			if result == 0 {
				err = l.svcCtx.Redis().Del(context.Background(), rediskey.UserConn(uid, platform, l.svcCtx.Config.PodIp())).Err()
				if err != nil {
					logx.Error("del user conn from redis err", " uid ", uid, " platform ", platform, " err ", err)
				}
			}
		}
	}()
	if oldStringMap, ok := l.WsConnToUser[conn]; ok {
		for k, v := range oldStringMap {
			platform = k
			uid = v
		}
		if oldConnMap, ok := l.WsUserToConn[uid]; ok {
			delete(oldConnMap, platform)
			l.WsUserToConn[uid] = oldConnMap
			if len(oldConnMap) == 0 {
				delete(l.WsUserToConn, uid)
			}
		}
		delete(l.WsConnToUser, conn)
	}
	err = conn.Close(websocket.StatusNormalClosure, "client close")
	if err != nil {
		logx.Error("close conn err", "", "uid", uid, "platform", platform)
		return err
	}
	return nil
}

func (l *Ws) GetUserUid(conn *UserConn) (uid string, platform string) {
	l.rwLock.RLock()
	defer l.rwLock.RUnlock()
	if stringMap, ok := l.WsConnToUser[conn]; ok {
		for k, v := range stringMap {
			platform = k
			uid = v
		}
		return uid, platform
	}
	return "", ""
}

func (l *Ws) checkUserConn(interval time.Duration) {
	l._checkUserConn("", "")
	//interval := time.Second * 60
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			l._checkUserConn("", "")
		}
	}
}

func (l *Ws) _checkUserConn(userId, pid string) {
	var uidPlatformMap = make(map[string][]string)
	l.rwLock.RLock()
	for uid, connMap := range l.WsUserToConn {
		for platform := range connMap {
			uidPlatformMap[uid] = append(uidPlatformMap[uid], platform)
		}
	}
	l.rwLock.RUnlock()
	for uid, platformList := range uidPlatformMap {
		for _, platformID := range platformList {
			if userId == "" || pid == "" || (pid == platformID && userId == uid) {
				err := l.svcCtx.Redis().HSet(
					context.Background(), rediskey.UserConn(uid, platformID, l.svcCtx.Config.PodIp()),
					fmt.Sprintf("%s%s", l.svcCtx.Config.PodIp(), l.svcCtx.Config.ListenOn), // key
					time.Now().UnixMilli(),                                                 // value
				).Err()
				if err != nil {
					logx.Error("add user conn to redis err ", " uid ", uid, " platform ", platformID, " err ", err)
					continue
				}
			}
		}
	}
}
