package nw

import (
	"babo/pb_code"

	"google.golang.org/protobuf/proto"
)

const (
	Tcp = "tcp"
	Ws  = "ws"
	Wss = "wss"
)

type ConnectFunc func(service IService, session ISession)
type DisConnectFunc func(service IService, session ISession)
type ReceiveFunc func(msg []byte) error

type IService interface {
	Start() error
	Stop() error
}

type ISession interface {
	ID() int32
	RemoteAddr() string

	Start(cb ReceiveFunc)
	ServeIO()
	Read()
	Write()

	SendData([]byte)
}

type IClient interface {
	SessionID() int32
	SendMsg(pb_code.MsgId, proto.Message)
	// GraceClose() error
	// Wait()
}
