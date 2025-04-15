package ws

import (
	"babo/utility/nw"
	"babo/utility/pkg"
	"babo/utility/zlog"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	DefaultReadBufferSize  int = 4096
	DefaultWriteBufferSize int = 4096
)

type Service struct {
	listener      net.Listener
	webscoket     websocket.Upgrader
	autoSessionID int32

	onConnect nw.ConnectFunc
	onDisConn nw.DisConnectFunc
}

type RemoteCtl struct {
	UseWhite  bool
	WhiteList map[string]bool
	UseBlack  bool
	BlackList map[string]bool
}

func (r *RemoteCtl) CheckOrigin(req *http.Request) bool {
	remoteAddr := req.RemoteAddr
	if index := strings.LastIndex(remoteAddr, ":"); index >= 0 {
		ip := remoteAddr[:index]

		if r.UseWhite {
			if _, ok := r.WhiteList[ip]; ok {
				return true
			}
			return false
		}

		if r.UseBlack {
			if _, ok := r.BlackList[ip]; ok {
				return false
			}
		}
		return true

	}
	zlog.Error("CheckOrigin, parse remote addr failed", zap.String("remoteAddr", remoteAddr))
	return true
}

func NewServive(ip string, port int32, ctl *RemoteCtl, onConnect nw.ConnectFunc, onDisConn nw.DisConnectFunc) (*Service, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zlog.Error("net.Listen error", zap.String("err", err.Error()))
		return nil, err
	}
	return &Service{
		listener: listener,
		webscoket: websocket.Upgrader{
			ReadBufferSize:    DefaultReadBufferSize,
			WriteBufferSize:   DefaultWriteBufferSize,
			WriteBufferPool:   &sync.Pool{},
			EnableCompression: false,
			CheckOrigin: func(r *http.Request) bool {
				return ctl.CheckOrigin(r)
			},
		},
		autoSessionID: 0,
		onConnect:     onConnect,
		onDisConn:     onDisConn,
	}, nil
}

func (s *Service) genSessionId() int32 {
	if s.autoSessionID == 0x7fffffff {
		s.autoSessionID = 0
	}
	return atomic.AddInt32(&s.autoSessionID, 1)
}

func (s *Service) Start() error {
	zlog.Info("WsService whill start")
	go func() {
		defer pkg.ProtectError()
		zlog.Info("WsService start", zap.String("addr", s.listener.Addr().String()))
		if err := http.Serve(s.listener, s); err != nil {
			zlog.Error("http.Serve error", zap.String("err", err.Error()))
			return
		}
	}()
	return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.webscoket.Upgrade(w, r, nil)
	if err != nil {
		zlog.Error("webscoket.Upgrade error", zap.String("err", err.Error()))
		return
	}
	session := NewSession(s.genSessionId(), conn)
	s.onConnect(s, session)
	go func() {
		defer pkg.ProtectError()
		session.ServeIO()
	}()
}

func (s *Service) Stop() error {
	return s.listener.Close()
}
