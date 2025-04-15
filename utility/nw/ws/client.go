package ws

import (
	"github.com/gorilla/websocket"

	"babo/utility/zlog"

	"go.uber.org/zap"
)

type Client struct {
	session *Session
}

func NewClient() *Client {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:10004", nil)
	if err != nil {
		zlog.Error("dial error", zap.Error(err))
		return nil
	}

	session := NewSession(1, conn)

	return &Client{
		session: session,
	}
}

func (c *Client) Init(session *Session) {
	c.session = session
	c.session.Start(c.onReceiveMsg)
}

func (c *Client) onReceiveMsg(data []byte) (err error) {
	// zlog.Info("onReceiveMsg", zap.String("data", string(data)))
	// rpc call gs

	return nil
}
