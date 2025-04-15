package handler

import (
	"babo/db"
	. "babo/game/common"
	"babo/game/match"
	"babo/game/player"
	"babo/game/room"
	"babo/orm"
	"babo/pb_code"
	"babo/utility/nw"
	pbhandler "babo/utility/pb_handler"
	"babo/utility/uuid"
	"babo/utility/zlog"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func init() {
	pbhandler.Register(pb_code.MsgId_Login, OnLoginReq)
	pbhandler.Register(pb_code.MsgId_HeartBeat, OnHeartBeatReq)
	pbhandler.Register(pb_code.MsgId_Match, OnMatchReq)
	pbhandler.Register(pb_code.MsgId_EnterRoom, OnEnterRoom)
}

func OnEnterRoom(c nw.IClient, r proto.Message) (proto.Message, error) {
	p := c.(*player.Player)
	msg := r.(*pb_code.Proto)
	req := &pb_code.EnterRoomReq{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		zlog.Error("unmarshal EnterRoomReq failed", zap.Error(err), ZapU(p))
		return nil, err
	}
	rsp := &pb_code.EnterRoomRsp{
		Code: pb_code.ResCode_Success,
	}

	room := room.Mgr.GetRoom(req.RoomId)
	if room == nil {
		zlog.Error("room not found", zap.Int64("room_id", req.RoomId), ZapU(p))
		rsp.Code = pb_code.ResCode_Fail
		return rsp, nil
	}

	if !room.IsValidPlayer(p) {
		zlog.Error("player not valid in room", zap.Int64("room_id", req.RoomId), ZapU(p))
		rsp.Code = pb_code.ResCode_Fail
		return rsp, nil
	}

	if room.IsPlayerEnter(p) {
		zlog.Error("player enter room repeatly", zap.Int64("room_id", req.RoomId), ZapU(p))
		rsp.Code = pb_code.ResCode_Fail
		return rsp, nil
	}

	rsp.Target = room.OnPlayerEnter(p)

	return rsp, nil
}

func OnMatchReq(c nw.IClient, r proto.Message) (proto.Message, error) {
	p := c.(*player.Player)
	msg := r.(*pb_code.Proto)
	req := &pb_code.MatchReq{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		zlog.Error("unmarshal MatchReq failed", zap.Error(err), ZapU(p))
		return nil, err
	}
	rsp := &pb_code.MatchRsp{
		Code: pb_code.ResCode_Success,
	}

	match.Mgr.PushMatch(p)

	return rsp, nil
}

func OnHeartBeatReq(c nw.IClient, r proto.Message) (proto.Message, error) {
	p := c.(*player.Player)
	msg := r.(*pb_code.Proto)
	req := &pb_code.HeartBeatReq{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		zlog.Error("unmarshal HeartBeatReq failed", zap.Error(err), ZapU(p))
		return nil, err
	}

	rsp := &pb_code.HeartBeatRsp{
		Time: time.Now().Unix(),
	}
	return rsp, nil
}

func OnLoginReq(c nw.IClient, r proto.Message) (proto.Message, error) {
	p := c.(*player.Player)
	msg := r.(*pb_code.Proto)
	req := &pb_code.LoginReq{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		zlog.Error("unmarshal LoginReq failed", zap.Error(err), ZapU(p))
		return nil, err
	}
	rsp := &pb_code.LoginRsp{
		Code: pb_code.ResCode_Success,
	}
	// todo sdk check

	// todo repeat login

	gorm_db := orm.DB()

	account := &db.AccountData{}
	account_res := gorm_db.Where("account = ?", req.Account).First(account)
	if account_res.Error != nil {
		if errors.Is(account_res.Error, gorm.ErrRecordNotFound) {
			// not found, create new account
			uid, err := uuid.Generate()
			if err != nil {
				zlog.Error("generate uuid failed", zap.Error(err), zap.String("account", req.Account))
				rsp.Code = pb_code.ResCode_Fail
				return nil, err
			}
			account.Account = req.Account
			account.Uid = uid
			if err := gorm_db.Create(account).Error; err != nil {
				zlog.Error("create account failed", zap.Error(err), zap.String("account", req.Account))
				rsp.Code = pb_code.ResCode_Fail
				return nil, err
			}
		} else {
			zlog.Error("query account failed", zap.Error(account_res.Error), zap.String("account", req.Account))
			return nil, account_res.Error
		}
	}

	user := &db.UserData{}
	user_res := gorm_db.Where("uid = ?", account.Uid).First(user)
	if user_res.Error != nil {
		if errors.Is(user_res.Error, gorm.ErrRecordNotFound) {
			// not found, create new user
			user.Uid = account.Uid
			user.Account = account.Account
			if err := gorm_db.Create(user).Error; err != nil {
				zlog.Error("create user failed", zap.Error(err), zap.String("account", req.Account))
				rsp.Code = pb_code.ResCode_Fail
				return nil, err
			}
		} else {
			zlog.Error("query user failed", zap.Error(user_res.Error), zap.String("account", req.Account))
			return nil, user_res.Error
		}
	}

	p.OnLogin(user)
	player.Mgr.StoreUid(p.SessionID(), user.Uid)

	rsp.Uid = user.Uid
	return rsp, nil
}
