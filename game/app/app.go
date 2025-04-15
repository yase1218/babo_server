package app

import (
	"babo/db"
	"babo/game/match"
	"babo/game/player"
	"babo/game/room"
	"babo/orm"
	"babo/utility/application"
	"babo/utility/nw"
	"babo/utility/nw/ws"
	"babo/utility/uuid"
	"babo/utility/zlog"
	"fmt"

	"go.uber.org/zap"
)

type Config struct {
	GameServer struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		WorkID       int64  `yaml:"workid"`
		DataCenterID int64  `yaml:"datacenterid"`
		JsonPath     string `yaml:"jsonpath"`
		Mysql        struct {
			Ip     string `yaml:"ip"`
			Port   string `yaml:"port"`
			User   string `yaml:"user"`
			Pwd    string `yaml:"pwd"`
			DBName string `yaml:"dbname"`
		} `yaml:"mysql"`
	} `yaml:"gameserver"`
}

type GameServerApp struct {
	*application.DefaultApplication
	config Config
	//gameRpcClient *rpc.RpcClient
}

func (sa *GameServerApp) GetConfig() *Config {
	return &sa.config
}

func (sa *GameServerApp) InitServices() error {
	// uuid
	if err := uuid.Init(sa.config.GameServer.WorkID,
		sa.config.GameServer.DataCenterID); err != nil {
		zlog.Error("init uuid failed", zap.String("err", err.Error()))
		return err
	}
	// player manager
	player.Mgr.Init()

	// room manager
	room.Mgr.Init()
	// match manager
	match.Mgr.Init()

	if err := orm.Connect(
		sa.config.GameServer.Mysql.Ip+":"+sa.config.GameServer.Mysql.Port,
		sa.config.GameServer.Mysql.User,
		sa.config.GameServer.Mysql.Pwd,
		sa.config.GameServer.Mysql.DBName); err != nil {
		zlog.Error("connect mysql failed", zap.String("err", err.Error()))
		return err
	}

	gorm_db := orm.DB()

	if gorm_db == nil {
		zlog.Error("gorm db is nil")
		return fmt.Errorf("gorm db is nil")
	}

	// init db
	if err := gorm_db.AutoMigrate(&db.AccountData{}); err != nil {
		zlog.Error("AutoMigrate AccountData failed", zap.String("err", err.Error()))
		return err
	}

	if err := gorm_db.AutoMigrate(&db.UserData{}); err != nil {
		zlog.Error("AutoMigrate UserData failed", zap.String("err", err.Error()))
		return err
	}

	// 初始化网络放在最后
	if err := sa.InitNet(nw.Ws); err != nil {
		return err
	}

	return nil
}

func (sa *GameServerApp) InitNet(netType string) error {
	// service for client
	switch netType {
	// case TcpNetService:
	// 	return NewTcpService(ip, port)
	case nw.Ws:
		cfg := sa.config.GameServer
		remotr_ctl := &ws.RemoteCtl{
			UseWhite:  false, //cfg.UseWhite,
			WhiteList: make(map[string]bool),
			UseBlack:  false, //cfg.UseBlack,
			BlackList: make(map[string]bool),
		}

		// for _, ip := range cfg.WhiteList {
		// 	remotr_ctl.WhiteList[ip] = true
		// }
		// for _, ip := range cfg.BlackList {
		// 	remotr_ctl.BlackList[ip] = true
		// }
		s, err := ws.NewServive(cfg.Host, int32(cfg.Port), remotr_ctl, sa.OnConnect, sa.OnDisconnect)
		if err != nil {
			return err
		}
		s.Start()
	// case WssService:
	// 	return NewWssService(ip, port, args[0].(*RemoteCtl))
	default:
		return fmt.Errorf("unknow service type: %s", netType)
	}

	// connect to gameserver
	// conn, _, err := websocket.DefaultDialer.Dial("ws://"+sa.config.GateServer.GameServer.Host+":"+strconv.Itoa(sa.config.GateServer.GameServer.Port), nil)
	// if err != nil {
	// 	zlog.Error("connect to game server failed", zap.String("err", err.Error()))
	// 	sa.playerService.Stop()
	// 	return err
	// }

	// gsSession := ws.NewSession(0, conn)
	// gsSession.Start(handler.ProcessGameMsg)
	// gs.GameConn.Init(gsSession)
	return nil
}

func (sa *GameServerApp) OnConnect(service nw.IService, session nw.ISession) {
	if p := player.Mgr.NewPlayer(session); p != nil {
		zlog.Info("New Client connected.", zap.String("remote addr", session.RemoteAddr()), zap.Int32("sid", session.ID()))
		p.Start()
	} else {
		zlog.Error("Create player failed", zap.String("remote addr", session.RemoteAddr()))
	}
}

func (sa *GameServerApp) OnDisconnect(service nw.IService, session nw.ISession) {
	zlog.Info("client disconnected", zap.String("remote addr", session.RemoteAddr()), zap.Int32("sid", session.ID())) // TODO print player info
	player.Mgr.DelPlayer(session.ID())
}
