package ws

import (
	"babo/utility/nw"
	"babo/utility/pkg"
	"babo/utility/zlog"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Session struct {
	id         int32
	remoteAddr string
	conn       *websocket.Conn
	closeWg    sync.WaitGroup
	onReceive  nw.ReceiveFunc
	sendQueue  chan []byte
	closeDone  sync.Once
}

func NewSession(id int32, conn *websocket.Conn) *Session {
	return &Session{
		id:         id,
		remoteAddr: conn.RemoteAddr().String(),
		conn:       conn,
		sendQueue:  make(chan []byte, 1024),
	}
}

func (s *Session) Start(cb nw.ReceiveFunc) {
	s.onReceive = cb
	if cb == nil {
		panic("recvCallback is nil")
	}

	go func() {
		defer pkg.ProtectError()
		s.ServeIO()
	}()
}

func (s *Session) ID() int32 {
	return s.id
}

func (s *Session) RemoteAddr() string {
	return s.remoteAddr
}

func (s *Session) ServeIO() {
	s.closeWg.Add(2)
	go s.Read()
	go s.Write()
}

func (s *Session) Read() {
	defer pkg.ProtectError()
	defer s.closeWg.Done()
	defer s.Close()
	for {
		msgType, data, err := s.conn.ReadMessage()
		if err != nil && err != io.EOF {
			//zlog.Error("ReadMessage error", zap.String("err", err.Error()))
			return
		}
		if msgType != websocket.BinaryMessage {
			zlog.Error("Invalid websocket msg type", zap.Int("type", msgType))
			return
		}
		if err = s.onReceive(data); err != nil {
			zlog.Error("onReceive error", zap.String("err", err.Error()))
			return
		}

	}
}

func (s *Session) Write() {
	defer pkg.ProtectError()
	defer s.closeWg.Done()
	defer s.Close()

	for v := range s.sendQueue {
		if v == nil {
			return
		}
		if s.conn == nil {
			return
		}
		writer, err := s.conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			zlog.Error("s.conn.NextWriter error", zap.String("err", err.Error()))
			return
		}

		for len(v) > 0 {
			n, err := writer.Write(v)
			if err != nil {
				zlog.Error("writer.Write error", zap.String("err", err.Error()))
				return
			}
			v = v[n:]
		}

		writer.Close()
	}
}

func (s *Session) Close() {
	s.closeDone.Do(func() {
		if s.conn != nil {
			s.sendQueue <- nil
			s.closeWg.Wait()
			s.conn.Close()
			s.conn = nil
		}
	})
}

func (s *Session) SendData(msg []byte) {
	s.sendQueue <- msg
}
