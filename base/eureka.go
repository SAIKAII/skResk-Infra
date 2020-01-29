package base

import (
	"time"

	"github.com/SAIKAII/skReskInfra"
	"github.com/kataras/iris"
	"github.com/tietang/go-eureka-client/eureka"
)

type EurekaStarter struct {
	infra.BaseStarter
	client *eureka.Client
}

func (e *EurekaStarter) Init(ctx infra.StarterContext) {
	e.client = eureka.NewClient(ctx.Props())
	e.client.Start()
}

func (e *EurekaStarter) Setup(ctx infra.StarterContext) {
	info := make(map[string]interface{})
	info["startTime"] = time.Now()
	info["appName"] = ctx.Props().GetDefault("app.name", "skResk")
	Iris().Get("/info", func(ctx iris.Context) {
		ctx.JSON(info)
	})
	Iris().Get("health", func(ctx iris.Context) {
		health := eureka.Health{
			Details: make(map[string]interface{}),
		}
		health.Status = eureka.StatusUp
		ctx.JSON(health)
	})
}
