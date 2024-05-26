package gmvc

import (
	"encoding/json"
	"net/http"
	"reflect"
	"sync"
)

// Validator 函数定义
type Validator func(ctx HttpContext, fieldMeta *ParamMeta, value interface{}) error

// Resolver 自定义解析器
type Resolver func(ctx HttpContext, fieldMeta *ParamMeta, origin string) (interface{}, error)

// HandleError error 全局处理器
type HandleError func(ctx HttpContext, err error) interface{}

// HandlerFunc 通用的HandlerFunc
type HandlerFunc func(ctx HttpContext)

// HandlerDispatcher 核心接口
type HandlerDispatcher interface {
	// RegisterValidator 注册Validator
	RegisterValidator(name string, v Validator)

	// RegisterResolver 注册Resolver
	RegisterResolver(name string, r Resolver)

	// RegisterResponsor 注册Responser
	RegisterResponsor(render RenderType, r Responsor)

	// SetErrorHandler 设置全局错误处理
	SetErrorHandler(eh HandleError)

	// InnerWrap 公共的wrap方法
	InnerWrap(handler Action) HandlerFunc
}

// DefaultDispatcher HandlerDispatcher缺省实现
type DefaultDispatcher struct {
	lock sync.RWMutex

	checkerMap  map[string]Validator
	resolverMap map[string]Resolver
	responsor   map[RenderType]Responsor

	errHandler HandleError
}

var _ HandlerDispatcher = &DefaultDispatcher{}

var (
	// nil
	nilValue *interface{}

	// DefaultInstance default instance
	DefaultInstance = &DefaultDispatcher{
		checkerMap:  make(map[string]Validator, 0),
		resolverMap: make(map[string]Resolver, 0),
		responsor:   make(map[RenderType]Responsor, 0),
		errHandler: func(ctx HttpContext, err error) interface{} {
			return err.Error()
		},
	}

	parser = HandlerParser{
		dispatcher: DefaultInstance,
	}

	// action handler 必须实现的接口
	handlerInterfaceType = reflect.TypeOf((*Action)(nil)).Elem()
)

func init() {
	// 注册默认渲染器
	DefaultInstance.RegisterResponsor(JSON, &JSONResponsor{})
	DefaultInstance.RegisterResponsor(HTML, &HTMLResponsor{})
	DefaultInstance.RegisterResponsor(String, &StringResponsor{})

	// 注册内置参数解析器
	DefaultInstance.RegisterResolver("Json", func(ctx HttpContext, fieldMeta *ParamMeta, origin string) (interface{}, error) {
		switch fieldMeta.fieldType.Type.Kind() {
		case reflect.Struct:
			jsonFieldValue := reflect.New(fieldMeta.fieldType.Type)
			_ = json.Unmarshal([]byte(origin), jsonFieldValue.Interface())

			return jsonFieldValue.Elem().Interface(), nil
		case reflect.Ptr:
			if fieldMeta.fieldType.Type.Elem().Kind() == reflect.Struct {
				jsonFieldValue := reflect.New(fieldMeta.fieldType.Type)
				_ = json.Unmarshal([]byte(origin), jsonFieldValue.Interface())
				return jsonFieldValue.Elem().Interface(), nil
			}
		}

		return nil, nil
	})
}

// RegisterValidator 注册全局参数验证器
func (instance *DefaultDispatcher) RegisterValidator(name string, c Validator) {
	instance.lock.Lock()
	defer instance.lock.Unlock()

	instance.checkerMap[name] = c
}

// RegisterResolver 注册全局参数验证器
func (instance *DefaultDispatcher) RegisterResolver(name string, r Resolver) {
	instance.lock.Lock()
	defer instance.lock.Unlock()

	instance.resolverMap[name] = r
}

// RegisterResponsor 注册返回器
func (instance *DefaultDispatcher) RegisterResponsor(render RenderType, r Responsor) {
	instance.lock.Lock()
	defer instance.lock.Unlock()

	instance.responsor[render] = r
}

// SetErrorHandler 设置全局错误处理
func (instance *DefaultDispatcher) SetErrorHandler(eh HandleError) {
	instance.lock.Lock()
	defer instance.lock.Unlock()

	instance.errHandler = eh
}

// InnerWrap 关键方法
func (instance *DefaultDispatcher) InnerWrap(handler Action) HandlerFunc {
	return instance.InnerWrap0(reflect.TypeOf(handler))
}

// InnerWrap0 关键方法实现
// 入参handler必须是一个func类型。
// 返回是是一个HandlerFunc
func (instance *DefaultDispatcher) InnerWrap0(handlerType reflect.Type) HandlerFunc {
	if !handlerType.Implements(handlerInterfaceType) {
		panic("handler must implement 'gmvc.Action' interface")
	}

	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
	}

	if handlerType.Kind() != reflect.Struct {
		panic("handler or *handler must be struct type")
	}

	handlerMeta := parser.introspect(handlerType)

	return func(ctx HttpContext) {
		defer onRecover(ctx)

		// 1. 解析参数
		handlerInstance, err := instance.resolve(ctx, handlerMeta)
		if err != nil {
			resp := instance.errHandler(ctx, err)
			instance.doResponse(ctx, resp)
			return
		}

		// 2. 调用Init方法
		err = instance.initialize(ctx, handlerInstance)
		if err != nil {
			resp := instance.errHandler(ctx, err)
			instance.doResponse(ctx, resp)
			return
		}

		// 3. 调用Go方法
		resp, err := instance.launch(ctx, handlerInstance)
		if err != nil {
			resp = instance.errHandler(ctx, err)
			instance.doResponse(ctx, resp)
			return
		}

		instance.doResponse(ctx, resp)
		return
	}
}

func (instance *DefaultDispatcher) initialize(ctx HttpContext, handler interface{}) error {
	if initer, ok := handler.(Initer); ok {
		return initer.Init()
	}

	return nil
}

func (instance *DefaultDispatcher) launch(ctx HttpContext, handler interface{}) (interface{}, error) {
	if entity, ok := handler.(Action); ok {
		return entity.Go()
	}

	// this shoud never happen.
	panic("handler is not implements 'gmvc.Action'.")
}

func (instance *DefaultDispatcher) doResponse(ctx HttpContext, resp interface{}) {
	if entity, ok := resp.(Response); ok {
		if responsor, ok := instance.responsor[entity.Render]; ok {
			responsor.Response(ctx, &entity)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	if entity, ok := resp.(*Response); ok {
		if responsor, ok := instance.responsor[entity.Render]; ok {
			responsor.Response(ctx, entity)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	instance.doResponse(ctx, &Response{
		Render:     JSON,
		Body:       resp,
		StatusCode: http.StatusOK,
	})
}

func onRecover(c HttpContext) {
	if x := recover(); x != nil {
		c.Status(http.StatusInternalServerError)
	}
}
