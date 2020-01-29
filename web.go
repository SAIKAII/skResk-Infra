package infra

var apiInitializerRegister *InitializerRegister = new(InitializerRegister)

// 注册Web API初始化对象
func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

// 获取注册的Web api初始化对象
func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
