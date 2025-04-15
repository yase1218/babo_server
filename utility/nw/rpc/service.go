package rpc

// import (
// 	"context"
// 	"errors"
// 	"net"

// 	"gitlab.wowstudio.com/server/rambo/utility/nw"
// 	"gitlab.wowstudio.com/server/rambo/utility/pkg"
// 	"gitlab.wowstudio.com/server/rambo/utility/zlog"

// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// 	"google.golang.org/grpc/metadata"
// )

// type RpcService struct {
// 	server *grpc.Server
// }

// func (rs *RpcService) Stop() {
// 	rs.server.Stop()
// }

// func (rs *RpcService) GracefulStop() {
// 	rs.server.GracefulStop()
// }

// func auth(ctx context.Context) error {
// 	//md, ok := metadata.FromIncomingContext(ctx)
// 	_, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		zlog.Error("get metadata failed")
// 		return errors.New("get metadata failed")
// 	}

// 	// var (
// 	// 	appid  string
// 	// 	appkey string
// 	// )

// 	// if v, ok := md["appid"]; ok {
// 	// 	appid = v[0]
// 	// }

// 	// if v, ok := md["appkey"]; ok {
// 	// 	appkey = v[0]
// 	// }

// 	// if appid != "10010" || appkey != "is a key" {
// 	// 	logins.Error("auth failed")
// 	// 	return errors.New("auth failed")
// 	// }
// 	return nil
// }

// func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	if err := auth(ctx); err != nil {
// 		return nil, err
// 	}
// 	return handler(ctx, req)
// }

// func NewService(addr string, OpenTLS bool, services map[*grpc.ServiceDesc]interface{}) (*RpcService, error) {
// 	listenr, err := net.Listen(nw.Tcp, addr)
// 	if err != nil {
// 		zlog.Error("listen failed")
// 		return nil, err
// 	}
// 	var opts []grpc.ServerOption

// 	if OpenTLS {
// 		// TODO: use config
// 		creds, cerr := credentials.NewServerTLSFromFile("../keys/rpcserver.pem", "../keys/rpcserver.key")
// 		if cerr != nil {
// 			zlog.Error("load tls failed", zap.String("err", cerr.Error()))
// 			return nil, cerr
// 		}
// 		opts = append(opts, grpc.Creds(creds))
// 	}

// 	//opts = append(opts, grpc.UnaryInterceptor(interceptor))
// 	s := grpc.NewServer(opts...)
// 	rs := &RpcService{
// 		server: s,
// 	}

// 	// for k, v := range services {
// 	// 	s.RegisterService(k, v)
// 	// }

// 	go func() {
// 		defer pkg.ProtectError()
// 		for k, v := range services {
// 			s.RegisterService(k, v)
// 		}
// 		zlog.Info("rpc server start", zap.String("addr", addr))
// 		if er := s.Serve(listenr); er != nil {
// 			zlog.Error("rpc server failed", zap.String("err", er.Error()))
// 			return
// 		}
// 		zlog.Info("rpc server stop")
// 	}()
// 	return rs, nil
// }
