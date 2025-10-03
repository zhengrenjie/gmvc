package gmvc

// Initer 参数解析完成之后执行的回调接口
type Initer interface {
	Init() error
}

// Action handler执行方法
type Action interface {
	Go() (interface{}, error)
}

// Validator 函数定义
type Validator func(ctx GmvcContext, fieldMeta *ParamMeta, value interface{}) error

// Resolver 自定义解析器
type Resolver func(ctx GmvcContext, fieldMeta *ParamMeta, origin string) (interface{}, error)

// HandleError error 全局处理器
type HandleError func(ctx GmvcContext, err error) interface{}

// HandlerFunc 通用的HandlerFunc
type HandlerFunc func(ctx GmvcContext)

// RecoverFunc panic之后gmvc会自动recover, 然后回调该接口进行recover逻辑执行
type RecoverFunc func(ctx GmvcContext, any interface{}) (interface{}, error)
