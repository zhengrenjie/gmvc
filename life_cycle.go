package gmvc

// Initer 参数解析完成之后执行的回调接口
type Initer interface {
	Init() error
}

// Action handler执行方法
type Action interface {
	Go() (interface{}, error)
}
