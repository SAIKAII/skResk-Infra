package lb

import (
	"fmt"
	"github.com/tietang/go-eureka-client/eureka"
	"strings"
)

var lbs = make(map[string]Balancer, 0)

type Apps struct {
	Client *eureka.Client
}

type App struct {
	Name      string
	Instances []*ServerInstance
	lb        Balancer
}

// AddBalancer 添加应用对应的负载均衡算法，或替换已有应用的负载均衡算法
func AddBalancer(appName string, lb Balancer) {
	lbs[appName] = lb
}

func (a *App) Get(key string) *ServerInstance {
	return a.lb.Next(key, a.Instances)
}

func (a *Apps) Get(appName string) *App {
	var application *eureka.Application
	for _, a := range a.Client.Applications.Applications {
		if a.Name == strings.ToUpper(appName) {
			application = &a
			break
		}
	}
	if application == nil {
		return nil
	}

	lb, ok := lbs[appName]
	if !ok {
		lb = &RoundRobinBalancer{}
		lbs[appName] = lb
	}

	na := &App{
		Name:      application.Name,
		Instances: make([]*ServerInstance, 0),
		lb:        lb,
	}
	for _, ins := range application.Instances {
		var port int
		if ins.SecurePort.Enabled {
			port = ins.SecurePort.Port
		} else {
			port = ins.Port.Port
		}
		si := &ServerInstance{
			InstanceId: ins.InstanceId,
			AppName:    ins.AppName,
			Address:    fmt.Sprintf("%s:%d", ins.IpAddr, port),
			Status:     Status(ins.Status),
			MetaData:   make(map[string]string),
		}
		si.MetaData["rpcAddr"] = fmt.Sprintf("%s:%s", ins.IpAddr, ins.Metadata.Map["rpcPort"])
		na.Instances = append(na.Instances, si)
	}
	return na
}

// 服务实例状态
type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
)

// 服务实例
type ServerInstance struct {
	InstanceId string
	AppName    string
	Address    string
	Status     Status
	MetaData   map[string]string
}
