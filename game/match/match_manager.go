package match

import (
	. "babo/game/common"
	"babo/game/player"
	"babo/game/room"
	"babo/utility/pkg"
	"babo/utility/uuid"
	"babo/utility/zlog"
	"time"

	"go.uber.org/zap"
)

type PreRoom struct {
	players []*player.Player
}

type MatchManager struct {
	match_list []*player.Player
	unit_chan  chan *player.Player
	ticker     *time.Ticker
}

var Mgr MatchManager

func (m *MatchManager) Init() {
	m.match_list = make([]*player.Player, 0)
	m.unit_chan = make(chan *player.Player, 1024)

	zlog.Info("MatchManager Start")
	m.ticker = time.NewTicker(time.Second)
	go m.serve_io()
}

func (m *MatchManager) Close() {
	// todo grace close
	zlog.Info("MatchManager Close")
}

func (m *MatchManager) PushMatch(p *player.Player) {
	select {
	case m.unit_chan <- p:
	default:
		if _, ok := <-m.unit_chan; !ok {
			zlog.Info("unit_chan closed", ZapU(p))
		} else {
			zlog.Info("unit_chan full", ZapU(p))
		}
	}
}

func (m *MatchManager) serve_io() {
	defer pkg.ProtectError()
	for {
		select {
		case p := <-m.unit_chan:
			m.add_to_match(p)
		case <-m.ticker.C:
			m.on_ticker(time.Now())
		}
	}
}

func (m *MatchManager) add_to_match(p *player.Player) {
	// zlog.Info("add_to_match", zap.Int32("session_id", unit.session_id), zap.Int64("uid", unit.uid))
	m.match_list = append(m.match_list, p)
}

func (m *MatchManager) on_ticker(now time.Time) {
	m.process_match()
}

func (m *MatchManager) process_match() {
	// todo process match
	// zlog.Info("process_match")
	if len(m.match_list) < 2 {
		return
	}

	zlog.Debug("process_match start",
		zap.Int("len", len(m.match_list)))

	room_list := make([]*PreRoom, 0)
	for len(m.match_list) >= 2 {
		pre_room := &PreRoom{
			players: []*player.Player{
				m.match_list[0],
				m.match_list[1],
			},
		}
		m.match_list = m.match_list[2:]
		room_list = append(room_list, pre_room)
	}

	zlog.Debug("process_match end",
		zap.Int("len", len(m.match_list)),
		zap.Int("pair_list len", len(room_list)))

	for _, pre_room := range room_list {
		room_id, err := uuid.Generate()
		if err != nil {
			zlog.Error("room_id_generator.Generate failed",
				zap.Any("pair", pre_room), zap.Error(err))
			continue
		}
		zlog.Debug("process_match room_id",
			zap.Any("pair", pre_room), zap.Int64("room_id", room_id))

		room.Mgr.NewRoom(room_id, pre_room.players)
	}
}
