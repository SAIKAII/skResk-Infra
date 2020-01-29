package infra

import (
	"github.com/tietang/props/kvs"
)

type BootApplication struct {
	conf           kvs.ConfigSource
	starterContext StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContext: StarterContext{},
	}
	b.starterContext[KeyProps] = conf
	return b
}

func (b *BootApplication) Start() {
	// 初始化starter
	b.init()
	// 安装starter
	b.setup()
	// 启动
	b.start()
}

func (b *BootApplication) init() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	for i, starter := range StarterRegister.AllStarters() {
		if starter.StartBlocking() {
			// 最后一个可阻塞的，直接执行
			if i+1 == len(StarterRegister.AllStarters()) {
				starter.Start(b.starterContext)
			} else {
				// 防止阻塞后面的starter
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)
		}
	}
}
