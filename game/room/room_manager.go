package room

import "babo/game/player"

type RoomManager struct {
	room_map map[int64]*Room
}

var Mgr RoomManager

func (r *RoomManager) Init() {
	r.room_map = make(map[int64]*Room)
}

func (r *RoomManager) NewRoom(room_id int64, players []*player.Player) *Room {
	room := &Room{
		room_id:    room_id,
		players:    make(map[int64]*player.Player),
		enter_uids: make(map[int64]interface{}),
	}
	for _, p := range players {
		room.players[p.Uid()] = p
	}
	r.room_map[room_id] = room
	room.NtfCreate()
	return room
}

func (r *RoomManager) GetRoom(room_id int64) *Room {
	if room, ok := r.room_map[room_id]; ok {
		return room
	}
	return nil
}
