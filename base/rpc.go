package base

import (
	"net"
	"net/rpc"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/SAIKAII/skReskInfra"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	logrus.Infof("goRPC Register: %s", typ.String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	s.server = rpc.NewServer()
	rpcServer = s.server
	//rpcServer.Register(&gorpc.EnvelopeRpc{})
}

func (s *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "8082")
	// 监听网络端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info("tcp port listened for rpc: ", port)
	// 处理网络连接和请求
	go s.server.Accept(listener)
}
