package player

import (
	"babo/utility/nw"
	"sync"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// PlayerManager
var Mgr PlayerManager

type PlayerManager struct {
	player_pool        sync.Pool
	player_session_map sync.Map
	player_id_map      sync.Map
}

func (p *PlayerManager) Init() {
	p.player_pool.New = func() interface{} {
		return &Player{
			session: nil,
		}
	}
}

func (p *PlayerManager) NewPlayer(session nw.ISession) *Player {
	player := p.player_pool.Get().(*Player)
	player.Init(session)
	p.player_session_map.Store(session.ID(), player)

	return player
}

func (p *PlayerManager) GetPlayer(sessionID int32) *Player {
	v, ok := p.player_session_map.Load(sessionID)
	if ok {
		return v.(*Player)
	}
	return nil
}

func (p *PlayerManager) DelPlayer(sessionID int32) {
	v, ok := p.player_session_map.Load(sessionID)
	if ok {
		// TODO 下线逻辑
		v.(*Player).UnInit()
		p.player_session_map.Delete(sessionID)
		p.player_pool.Put(v)
	}
}

func (p *PlayerManager) StoreUid(sessionID int32, uid int64) {
	v, ok := p.player_session_map.Load(sessionID)
	if ok {
		p.player_id_map.Store(uid, v.(*Player))
	}
}

func (p *PlayerManager) GetPlayerByUid(uid int64) *Player {
	v, ok := p.player_id_map.Load(uid)
	if ok {
		return v.(*Player)
	}
	return nil
}
