package player

import (
	"babo/db"
	. "babo/game/common"
	"babo/pb_code"
	"babo/utility/nw"
	pbhandler "babo/utility/pb_handler"
	"babo/utility/pkg"
	"babo/utility/zlog"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type OutMsg struct {
	Id  int32
	Msg proto.Message
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////
// Player

type Player struct {
	session   nw.ISession
	user_data *db.UserData

	recv_chan_client chan []byte
	//route_msg_chan   chan *OutMsg

	ticker      *time.Ticker
	logout_chan chan struct{}
	close_chan  chan struct{}
	close_once  sync.Once
}

func (p *Player) GetSession() nw.ISession {
	return p.session
}

func (p *Player) Account() string {
	if p.user_data != nil {
		return p.user_data.Account
	}
	return ""
}

func (p *Player) Uid() int64 {
	if p.user_data != nil {
		return p.user_data.Uid
	}
	return 0
}

func (p *Player) SessionID() int32 {
	if p.session == nil {
		return 0
	}
	return p.session.ID()
}

func (p *Player) Init(session nw.ISession) {
	p.session = session

	p.recv_chan_client = make(chan []byte, 1024)
	//p.route_msg_chan = make(chan *OutMsg, 1024)
	p.close_chan = make(chan struct{})
	p.logout_chan = make(chan struct{}, 1)
	p.close_once = sync.Once{}
}

func (p *Player) UnInit() {
	p.session = nil
}

func (p *Player) Str() string {
	if p.session != nil {
		return fmt.Sprintf("[A:%s, U:%d, S:%d, ptr:%p]", p.Account(), p.Uid(), p.SessionID(), p)
	} else {
		return fmt.Sprintf("[A:%s, U:%d, ptr:%p]", p.Account(), p.Uid(), p)
	}
}

func (p *Player) Start() {
	p.ticker = time.NewTicker(time.Second)
	go p.serve_io()
	p.session.Start(p.OnRecv)
}

func (p *Player) OnRecv(m []byte) error {
	select {
	case p.recv_chan_client <- m:
	default:
		if _, ok := <-p.recv_chan_client; !ok {
			zlog.Info("recv_chan_client closed", ZapU(p))
		} else {
			zlog.Info("recv_chan_client full", ZapU(p))
		}
	}

	return nil
}

func (p *Player) SendMsg(id pb_code.MsgId, m proto.Message) {
	proto_msg := &pb_code.Proto{
		Id: id,
	}

	if m != nil {
		data, err := proto.Marshal(m)
		if err != nil {
			zlog.Error("SendMsg base proto.Marshal error", zap.Error(err), ZapU(p))
			return
		}
		proto_msg.Body = data
	}
	buff, err := proto.Marshal(proto_msg)
	if err != nil {
		zlog.Error("SendMsg normal proto.Marshal error", zap.Error(err), ZapU(p))
		return
	}

	p.session.SendData(buff)
}

func (p *Player) serve_io() {
	defer pkg.ProtectError()
	for {
		select {
		case m := <-p.recv_chan_client:
			p.handle_client_msg(m)
		// case m := <-p.recv_chan_s:
		// 	p.handle_social_msg(m)
		// case order := <-p.order_chan:
		// 	p.process_order(order)
		case <-p.ticker.C:
			p.on_ticker(time.Now())
		case <-p.logout_chan:
			p.logout()
		// case m := <-p.route_msg_chan:
		// 	p.SendMsg(m.Id, m.Msg)
		case <-p.close_chan:
			p.ticker.Stop()
			close(p.recv_chan_client)
			// close(p.recv_chan_s)
			// close(p.order_chan)
			close(p.logout_chan)
			//close(p.route_msg_chan)
			return
		}
	}
}

func (p *Player) OnLogout() {
	select {
	case p.logout_chan <- struct{}{}:
	default:
		if _, ok := <-p.logout_chan; !ok {
			zlog.Info("logout_chan closed", ZapU(p))
		} else {
			zlog.Info("logout_chan full", ZapU(p))
		}
	}
}

func (p *Player) on_ticker(now time.Time) {
}

func (p *Player) handle_client_msg(m []byte) {
	zlog.Debug("handle_client_msg", zap.String("msg", string(m)), ZapU(p))

	proto_msg := &pb_code.Proto{}
	err := proto.Unmarshal(m, proto_msg)
	if err != nil {
		zlog.Error("handle_client_msg proto.Unmarshal error", zap.Error(err))
		return
	}
	zlog.Debug("handle_client_msg proto_msg", zap.Any("proto_msg", proto_msg), ZapU(p))

	handler := pbhandler.GetHandler(proto_msg.Id)
	if handler == nil {
		zlog.Error("handler is nil", zap.Int32("id", int32(proto_msg.Id)), ZapU(p))
		return
	}

	rsp, err := handler(p, proto_msg)
	if err != nil {
		zlog.Error("handle_client_msg handler error", zap.Error(err), ZapU(p))
		return
	}

	if rsp == nil {
		zlog.Error("handle_client_msg handler rsp", zap.Any("rsp", rsp), ZapU(p))
		return
	}

	p.SendMsg(proto_msg.Id, rsp)
}

func (p *Player) logout() {
	defer pkg.ProtectError()

	p.close_once.Do(func() {
		// TODO 处理下线逻辑

		zlog.Info("Player Logout", ZapU(p))

		close(p.close_chan)

		Mgr.DelPlayer(p.session.ID())
		p.UnInit()
	})
}

func (p *Player) OnLogin(user_data *db.UserData) {
	p.user_data = user_data
}
