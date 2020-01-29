package infra

import (
	"sort"

	"github.com/tietang/props/kvs"
)

const (
	KeyProps = "_conf"
)

// 基础资源上下文结构提
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没被初始化")
	}
	return p.(kvs.ConfigSource)
}

func (s StarterContext) SetProps(conf kvs.ConfigSource) {
	s[KeyProps] = conf
}

type PriorityGroup int

const (
	SystemGroup         PriorityGroup = 30
	BasicResourcesGroup PriorityGroup = 20
	AppGroup            PriorityGroup = 10

	INT_MAX          = int(^uint(0) >> 1)
	DEFAULT_PRIORITY = 10000
)

// 基础资源启动器接口
type Starter interface {
	// 系统启动，初始化一些基础资源
	Init(StarterContext)
	// 系统基础资源的安装
	Setup(StarterContext)
	// 启动基础资源
	Start(StarterContext)
	// 启动器是否可以阻塞
	StartBlocking() bool
	// 资源停止和销毁
	Stop(StarterContext)
	PriorityGroup() PriorityGroup
	Priority() int
}

// 基础空启动器实现，为了方便资源启动器的代码实现
type BaseStarter struct {
}

func (b *BaseStarter) Init(StarterContext)          {}
func (b *BaseStarter) Setup(StarterContext)         {}
func (b *BaseStarter) Start(StarterContext)         {}
func (b *BaseStarter) StartBlocking() bool          { return false }
func (b *BaseStarter) Stop(StarterContext)          {}
func (b *BaseStarter) PriorityGroup() PriorityGroup { return BasicResourcesGroup }
func (b *BaseStarter) Priority() int                { return DEFAULT_PRIORITY }

// 启动器注册器
type starterRegister struct {
	// nonBlockingStarters []Starter
	// blockingStarters    []Starter
	starters []Starter
}

func (r *starterRegister) Register(s Starter) {
	// if s.StartBlocking() {
	// 	r.blockingStarters = append(r.blockingStarters, s)
	// } else {
	// 	r.nonBlockingStarters = append(r.nonBlockingStarters, s)
	// }
	// typ := reflect.TypeOf(s)
	// log.Printf("Register starter: %s", typ.String())
	r.starters = append(r.starters, s)
}

func (r *starterRegister) AllStarters() []Starter {
	// starters := make([]Starter, 0)
	// starters = append(starters, r.nonBlockingStarters...)
	// starters = append(starters, r.blockingStarters...)
	// return starters
	return r.starters
}

var StarterRegister *starterRegister = &starterRegister{}

type Starters []Starter

func (s Starters) Len() int {
	return len(s)
}

func (s Starters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}

func Register(s Starter) {
	StarterRegister.Register(s)
}

func SortStarters() {
	sort.Sort(Starters(StarterRegister.AllStarters()))
}

func GetStarters() []Starter {
	return StarterRegister.AllStarters()
}
