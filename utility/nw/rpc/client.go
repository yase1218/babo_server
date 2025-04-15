package rpc

// import (
// 	"context"

// 	"gitlab.wowstudio.com/server/rambo/utility/zlog"

// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type RpcClient struct {
// 	conn *grpc.ClientConn
// }

// // client := pbgo.NewHelloClient(conn)
// //
// //	for {
// //		zlog.Info("request", zap.String("msg", "fucker"))
// //		resp, er := client.SayHello(context.Background(), &pbgo.HelloRequest{Name: "fucker"})
// //		if er != nil {
// //			zlog.Error("rpc call failed", zap.String("er", er.Error()))
// //			return
// //		}
// //		zlog.Info("resp", zap.String("msg", resp.GetMessage()))
// //		time.Sleep(3 * time.Second)
// //	}
// func (c *RpcClient) Conn() *grpc.ClientConn {
// 	return c.conn
// }

// func (c *RpcClient) Close() {
// 	c.conn.Close()
// }

// // customCredential is an implementation of credentials.PerRPCCredentials.
// type customCredential struct{}

// // GetRequestMetadata gets the current request metadata.
// func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
// 	// TODO: use config
// 	return map[string]string{
// 		"appid":  "10010",
// 		"appkey": "is a key",
// 	}, nil
// }

// // RequireTransportSecurity returns true if the credentials requires transport security.
// func (c customCredential) RequireTransportSecurity() bool {
// 	return false // TODO: 暂时关闭自定义认证
// }

// func NewClient(addr string, openTLS bool) (*RpcClient, error) {
// 	var opts []grpc.DialOption
// 	if openTLS {
// 		// TODO: use config
// 		creds, err := credentials.NewClientTLSFromFile("../keys/rpcserver.pem", "*.wowstudio.com")
// 		if err != nil {
// 			zlog.Error("load tls failed", zap.String("err", err.Error()))
// 			return nil, err
// 		}
// 		opts = append(opts, grpc.WithTransportCredentials(creds))
// 	} else {
// 		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	}

// 	//opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
// 	// conn, er := grpc.Dial(Addr, grpc.WithTransportCredentials(creds))
// 	zlog.Info("rpc connect to", zap.String("addr", addr))
// 	conn, er := grpc.Dial(addr, opts...)
// 	if er != nil {
// 		zlog.Error("connect failed", zap.String("err", er.Error()))
// 		return nil, er
// 	}
// 	// defer conn.Close()

// 	return &RpcClient{conn: conn}, nil
// }
