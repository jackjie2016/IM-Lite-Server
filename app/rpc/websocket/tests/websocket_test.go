package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/websocketservice"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"strconv"
	"strings"
	"testing"
	"time"
)

var url = "ws://localhost:31000?userID=%s&platform=%s&token=%s"

func TestWebsocket1(t *testing.T) {
	var (
		userId   = "1"
		platform = "Mobile"
		token    = "test-token"
	)
	conn, response, err := websocket.Dial(ctx, fmt.Sprintf(url, userId, platform, token), nil)
	if err != nil {
		t.Fatalf("Dial error: %v\nresponse: %+v", err, response)
	}
	go testKeepAlive(userId, t)
	readData(conn, t)
}

func testKeepAlive(uid string, t *testing.T) {
	f := func() {
		/*_, err := websocketService.ListUserConversation(ctx, &pb.ListUserConversationReq{UserId: uid})
		if err != nil {
			t.Fatalf("ListUserConversation error: %v", err)
			return
		}*/
		request, _ := http.NewRequest("POST", "http://localhost:8888/websocket/v1/conversation/list", nil)
		request.Header.Set("token", "token.uid."+uid)
		_, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatalf("ListUserConversation error: %v", err)
			return
		}
	}
	f()
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}

func readData(conn *websocket.Conn, t *testing.T) {
	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("read err:", err)
			time.Sleep(time.Second)
			continue
		}
		body := &pb.PushBody{}
		e := proto.Unmarshal(message, body)
		if e == nil {
			msgData := &pb.MsgData{}
			e = proto.Unmarshal(body.Data, msgData)
			if e == nil {
				t.Logf("msgData: %+v", msgData.String())
			}
		}
	}
}

func TestWebsocketCreateRoom1(t *testing.T) {
	conversation, err := websocketService.CreateConversation(ctx, &pb.CreateConversationReq{Members: []string{"unit-test-0", "unit-test-1", "1"}})
	if err != nil {
		t.Fatalf("CreateConversation error: %v", err)
		return
	}
	fmt.Println(conversation.String())
	// 62d2375770474c14c0a09906
}

func TestWebsocketSendMsg1(t *testing.T) {
	msg, err := websocketService.SendMsg(ctx, &pb.SendMsgReq{Msg: &pb.MsgData{
		ClientMsgID: "4",
		SenderID:    "unit-test-0",
		ConvID:      "62d2375770474c14c0a09906",
		ContentType: 0,
		Content:     []byte(fmt.Sprintf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))),
		ClientTime:  time.Now().UnixMilli(),
		ServerTime:  0,
		Seq:         0,
		OfflinePush: nil,
		MsgOptions: &pb.MsgOptions{
			Storage: true,
			Unread:  true,
		},
	}})
	if err != nil {
		t.Fatalf("SendMsg error: %v", err)
		return
	}
	fmt.Println(msg.String())
}
func TestWebsocketSendMsg11(t *testing.T) {
	msgData := &pb.MsgData{
		ClientMsgID: "4",
		SenderID:    "unit-test-0",
		ConvID:      "62d2375770474c14c0a09906",
		ContentType: 0,
		Content:     []byte(fmt.Sprintf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))),
		ClientTime:  time.Now().UnixMilli(),
		ServerTime:  0,
		Seq:         0,
		OfflinePush: nil,
		MsgOptions: &pb.MsgOptions{
			Storage: true,
			Unread:  true,
		},
	}
	buf, _ := proto.Marshal(msgData)
	req, _ := http.NewRequest("POST", "http://localhost:8888/websocket/v1/message/send", bytes.NewReader(buf))
	req.Header.Set("token", "token.uid.unit-test-0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("SendMsg error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestWebsocketSendMsg2(t *testing.T) {
	const num = 3
	var fss [num][]func()
	for i := 0; i < 1000; i++ {
		fss[i%num] = append(fss[i%num], func() {
			msg, err := websocketService.SendMsg(ctx, &pb.SendMsgReq{Msg: &pb.MsgData{
				ClientMsgID: strconv.FormatInt(time.Now().UnixNano(), 10),
				SenderID:    "unit-test-0",
				ConvID:      "62cfffbc8fd5380d0202985e",
				ContentType: 0,
				Content:     []byte(fmt.Sprintf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))),
				ClientTime:  time.Now().UnixMilli(),
				ServerTime:  0,
				Seq:         0,
				OfflinePush: nil,
				MsgOptions: &pb.MsgOptions{
					Storage: true,
					Unread:  true,
				},
			}})
			if err != nil {
				t.Fatalf("SendMsg error: %v", err)
				return
			}
			fmt.Println(msg.String())
		})
	}
	var fs []func() error
	for _, fss := range fss {
		fs = append(fs, func() error {
			for _, f := range fss {
				f()
			}
			return nil
		})
	}
	start := time.Now()
	mr.Finish(fs...)
	fmt.Println("send msg:", time.Now().Sub(start).Seconds())
}

/************************************************************************************************************/
func loadByTest(serviceName string, cfg interface{}) error {
	err := conf.Load(fmt.Sprintf(`../etc/%s.yaml`, serviceName), cfg)
	return err
}

type rpcCfg struct {
	ListenOn string
}

func (r *rpcCfg) Endpoints() []string {
	if strings.HasPrefix(r.ListenOn, ":") {
		return []string{"127.0.0.1" + r.ListenOn}
	}
	return []string{r.ListenOn}
}

var (
	ctx              = context.Background()
	websocketService websocketservice.WebsocketService
)

func init() {
	initWebsocket()
}

func initWebsocket() {
	cfg := &rpcCfg{}
	err := loadByTest("websocket", cfg)
	if err == nil {
		client, err := zrpc.NewClient(zrpc.RpcClientConf{
			Endpoints: cfg.Endpoints(),
		})
		if err == nil {
			websocketService = websocketservice.NewWebsocketService(client)
		}
	}
}
