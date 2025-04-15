package main

import (
	"babo/pb_code"
	"babo/utility/pkg"
	"babo/utility/zlog"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) Start() {
	// 连接服务器
	var err error
	c.Conn, _, err = websocket.DefaultDialer.Dial("ws://8.153.203.242:10005", nil)
	if err != nil {
		panic(err)
	}
	defer c.Conn.Close()

	login_req := &pb_code.LoginReq{
		Account: "test",
	}

	c.SendMsg(int32(pb_code.MsgId_Login), login_req)

	// 接收消息
	err = c.ReceiveMsg()
	if err != nil {
		panic(err)
	}
}

func (c *Client) SendMsg(msg_id int32, m proto.Message) {
	buff, err := proto.Marshal(m)
	if err != nil {
		zlog.Error("SendMsg proto.Marshal error", zap.Error(err))
		return
	}
	proto_msg := &pb_code.Proto{
		Id:   pb_code.MsgId(msg_id),
		Body: buff,
	}

	send_buff, err := proto.Marshal(proto_msg)
	if err != nil {
		zlog.Error("SendMsg proto.Marshal error", zap.Error(err))
		return
	}

	err = c.Conn.WriteMessage(websocket.BinaryMessage, send_buff)
	if err != nil {
		zlog.Error("SendMsg WriteMessage error", zap.Error(err))
		return
	}
	zlog.Info("SendMsg", zap.Any("proto_msg", proto_msg))
}

func (c *Client) ReceiveMsg() error {
	defer pkg.ProtectError()

	// 从服务器读取消息
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			zlog.Error("ReceiveMsg ERROR: ", zap.Error(err))
			return err
		}
		proto_msg := &pb_code.Proto{}
		err = proto.Unmarshal(message, proto_msg)
		if err != nil {
			zlog.Error("ReceiveMsg proto.Unmarshal error", zap.Error(err))
			return err
		}

		switch proto_msg.Id {
		case pb_code.MsgId_Login:
			login_rsp := &pb_code.LoginRsp{}
			err = proto.Unmarshal(proto_msg.Body, login_rsp)
			if err != nil {
				zlog.Error("ReceiveMsg proto.Unmarshal error", zap.Error(err))
				return err
			}
			zlog.Info("ReceiveMsg LoginRsp", zap.Any("rsp", login_rsp))

		default:
			zlog.Error("ReceiveMsg unknown msg id", zap.Any("id", proto_msg.Id))
		}
	}

	return nil
}
