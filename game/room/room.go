package room

import (
	"babo/game/player"
	"babo/pb_code"
)

type Room struct {
	room_id    int64
	players    map[int64]*player.Player
	enter_uids map[int64]interface{}
}

func (r *Room) NtfCreate() {
	// TODO 广播房间创建
	for _, p := range r.players {
		tar := r.GetTarget(p)
		if tar == nil {
			ntf := &pb_code.MatchResultNtf{
				RoomId: r.room_id,
			}

			p.SendMsg(pb_code.MsgId_MatchResult, ntf)
		} else {
			ntf := &pb_code.MatchResultNtf{
				RoomId: r.room_id,
				Target: &pb_code.MatchTarget{
					Uid: tar.Uid(),
				},
			}
			p.SendMsg(pb_code.MsgId_MatchResult, ntf)
		}
	}
}

func (r *Room) NtfEnter(p *player.Player) {
	// TODO 广播房间进入
	for _, v := range r.players {
		if p.Uid() != v.Uid() {
			ntf := &pb_code.UserEnterRoomNtf{
				Target: &pb_code.RoomTarget{
					Uid: p.Uid(),
				},
			}
			v.SendMsg(pb_code.MsgId_UserEnterRoom, ntf)
		}
	}
}

func (r *Room) GetTarget(p *player.Player) *player.Player {
	for _, v := range r.players {
		if p.Uid() != v.Uid() {
			return v
		}
	}
	return nil
}

func (r *Room) IsValidPlayer(p *player.Player) bool {
	_, ok := r.players[p.Uid()]
	return ok
}

func (r *Room) IsPlayerEnter(p *player.Player) bool {
	_, ok := r.enter_uids[p.Uid()]
	return ok
}

func (r *Room) OnPlayerEnter(p *player.Player) (res *pb_code.RoomTarget) {
	if len(r.enter_uids) > 0 {
		res = &pb_code.RoomTarget{}
		for uid := range r.enter_uids {
			res.Uid = uid
		}
	}

	r.NtfEnter(p)
	r.enter_uids[p.Uid()] = nil
	return
}
