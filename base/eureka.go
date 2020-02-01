package base

import (
	"time"

	"github.com/SAIKAII/skResk-Infra"
	"github.com/kataras/iris"
	"github.com/tietang/go-eureka-client/eureka"
)

var eurekaClient *eureka.Client

type EurekaStarter struct {
	infra.BaseStarter
	client *eureka.Client
}

func (e *EurekaStarter) Init(ctx infra.StarterContext) {
	e.client = eureka.NewClient(ctx.Props())
	rpcPort := ctx.Props().GetDefault("application.rpc.port", "8082")
	e.client.InstanceInfo.Metadata.Map["rpcPort"] = rpcPort
	e.client.Start()
	e.client.Applications, _ = e.client.GetApplications()
	eurekaClient = e.client
}

func (e *EurekaStarter) Setup(ctx infra.StarterContext) {
	info := make(map[string]interface{})
	info["startTime"] = time.Now()
	info["appName"] = ctx.Props().GetDefault("app.name", "skResk")
	Iris().Get("/info", func(ctx iris.Context) {
		ctx.JSON(info)
	})
	Iris().Get("/health", func(ctx iris.Context) {
		health := eureka.Health{
			Details: make(map[string]interface{}),
		}
		health.Status = eureka.StatusUp
		ctx.JSON(health)
	})
}
