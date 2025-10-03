package gmvc

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

var (
	// action handler
	handlerInterfaceType = reflect.TypeOf((*Action)(nil)).Elem()
)

func CreateGmvcBuilder(options ...GmvcOption) *GmvcBuilder {
	builder := &GmvcBuilder{
		actions:       make(map[string]HandlerFunc),
		checkerMap:    make(map[string]Validator),
		resolverMap:   make(map[string]Resolver),
		typedResolver: make(map[reflect.Type]Resolver),
		responsor:     make(map[RenderType]Responsor),
		globalMidware: make([]IMiddleware, 0),
		errHandler: func(ctx GmvcContext, err error) interface{} {
			return err.Error()
		},
		options: GmvcOptions{
			autodef: AnySrc, // 默认为任意位置
		},
		singletons: &singletonContext{
			typemap: make(map[reflect.Type]*singleton),
			namemap: make(map[string]*singleton),
		},
	}

	for _, option := range options {
		option(builder.options)
	}

	// Register default response render
	builder.RegisterResponsor(JSON, &JSONResponsor{})
	builder.RegisterResponsor(HTML, &HTMLResponsor{})
	builder.RegisterResponsor(String, &StringResponsor{})

	// Register default resolver
	builder.RegisterResolver("Json", func(ctx GmvcContext, fieldMeta *ParamMeta, origin string) (interface{}, error) {
		switch fieldMeta.fieldType.Type.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice:
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

	return builder
}

// var _ IGmvcRoute = (*route)(nil)

// type route struct {
// 	methods []string
// 	path    string
// 	action  HandlerFunc
// }

// // Action implements IGmvcRoute.
// func (r *route) Action() HandlerFunc {
// 	return r.action
// }

// // Method implements IGmvcRoute.
// func (r *route) Methods() []string {
// 	return r.methods
// }

// // Path implements IGmvcRoute.
// func (r *route) Path() string {
// 	return r.path
// }

// GmvcBuilder 负责创建Gmvc实例
type GmvcBuilder struct {
	actions       map[string]HandlerFunc
	checkerMap    map[string]Validator
	resolverMap   map[string]Resolver
	typedResolver map[reflect.Type]Resolver
	responsor     map[RenderType]Responsor
	errHandler    HandleError
	recover       RecoverFunc

	// 注册进gmvc的全局Middleware
	// 先注册先执行
	globalMidware []IMiddleware

	singletons *singletonContext

	options GmvcOptions
}

// RegisterRecover 注册Recover回调
func (gmvc *GmvcBuilder) RegisterRecover(recover RecoverFunc) *GmvcBuilder {
	gmvc.recover = recover
	return gmvc
}

// RegisterValidator 注册全局参数验证器
func (gmvc *GmvcBuilder) RegisterValidator(name string, c Validator) *GmvcBuilder {
	gmvc.checkerMap[name] = c
	return gmvc
}

// RegisterResolver 注册全局参数验证器
func (gmvc *GmvcBuilder) RegisterResolver(name string, r Resolver) *GmvcBuilder {
	gmvc.resolverMap[name] = r
	return gmvc
}

// RegisterTypedResolver 注册全局参数验证器
func (gmvc *GmvcBuilder) RegisterTypedResolver(typ []reflect.Type, r Resolver) *GmvcBuilder {
	for _, typ0 := range typ {
		gmvc.typedResolver[typ0] = r
	}

	return gmvc
}

// RegisterSingleton 注册需要组装到action中的实例
func (gmvc *GmvcBuilder) RegisterSingleton(name string, obj interface{}) *GmvcBuilder {
	return gmvc.registerSingleton(name, obj, false)
}

// RegisterSingleton 注册需要组装到action中的实例
func (gmvc *GmvcBuilder) registerSingleton(name string, obj interface{}, omitdup bool) *GmvcBuilder {
	typ := reflect.TypeOf(obj)
	if len(name) == 0 {
		name = typ.Name()
	}

	singleton := singleton{
		name: name,
		typ:  typ,
		obj:  obj,
	}

	if _, ok := gmvc.singletons.namemap[name]; ok {
		if !omitdup {
			panic("you are trying to register a exist name")
		}
	}

	gmvc.singletons.typemap[typ] = &singleton
	gmvc.singletons.namemap[name] = &singleton
	return gmvc
}

// RegisterResponsor 注册返回器
func (gmvc *GmvcBuilder) RegisterResponsor(render RenderType, r Responsor) *GmvcBuilder {
	gmvc.responsor[render] = r
	return gmvc
}

// RegisterResponsor 注册返回器
func (gmvc *GmvcBuilder) AddMiddleware(midware IMiddleware) *GmvcBuilder {
	gmvc.globalMidware = append(gmvc.globalMidware, midware)
	return gmvc
}

// SetErrorHandler 设置全局错误处理
func (gmvc *GmvcBuilder) SetErrorHandler(eh HandleError) *GmvcBuilder {
	gmvc.errHandler = eh
	return gmvc
}

type IActionBuilder interface {
	BuildAction(action Action, midware ...IMiddleware) HandlerFunc
}

func (gmvc *GmvcBuilder) ActionBuilder() IActionBuilder {
	return gmvc
}

func (gmvc *GmvcBuilder) BuildAction(action Action, midware ...IMiddleware) HandlerFunc {
	actionValue := reflect.ValueOf(action)
	handlerType := reflect.TypeOf(action)

	// action 必须是Action类型
	if !handlerType.Implements(handlerInterfaceType) {
		panic("handler must implement 'gmvc.Action' interface")
	}

	// 如果是指针，则取脂针所指的类型
	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
		actionValue = actionValue.Elem()
	}

	// 必须为一个struct类型
	if handlerType.Kind() != reflect.Struct {
		panic("handler or *handler must be struct type")
	}

	// 内省，获取action所有的元数据
	actionMeta := gmvc.introspect(actionValue)

	// 解析autowire然后注册进上下文
	autowires := actionMeta.GetAutowireInstances()
	for autowire, instance := range autowires {
		gmvc.registerSingleton(autowire, instance, true)
	}

	// 组装middleware
	// 最内层的执行方法，执行action的具体逻辑
	actionFunc := func(ctx GmvcContext) (interface{}, error) {
		// 1. 解析参数
		handlerInstance, err := gmvc.resolve(ctx, actionMeta)
		if err != nil {
			return nil, err
		}

		// 2. 调用Init方法
		err = gmvc.initialize(ctx, handlerInstance)
		if err != nil {
			return nil, err
		}

		// 3. 调用Go方法
		resp, err := gmvc.launch(ctx, handlerInstance)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	// 封装执行middleware的方法
	// 最内层的next方法就是action
	var next Next = actionFunc
	midwares := make([]IMiddleware, 0, len(gmvc.globalMidware)+len(midware))
	midwares = append(midwares, gmvc.globalMidware...)
	midwares = append(midwares, midware...)
	for i := len(midwares) - 1; i >= 0; i-- {
		midware := midwares[i]
		curnext := next
		var next0 Next = func(ctx GmvcContext) (interface{}, error) {
			// 判断当前middleware是否需要执行
			// 不需要执行，直接执行next
			if !midware.IsApply(ctx) {
				return curnext(ctx)
			}

			// 首先执行Before逻辑，请求可以被短路，只要返回任意结果就会提前结束，剩下的middleware和action都不会继续执行
			ret, err := midware.Before(ctx)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}

			// 执行around逻辑
			ret, err = midware.Around(ctx, curnext)
			// if err != nil {
			// 	return nil, err
			// }

			// if ret != nil {
			// 	return ret, nil
			// }

			// 最后执行after逻辑
			return midware.After(ctx, ret, err)
		}

		next = next0
	}

	// the outer handlerfunc
	handlerfunc := func(ctx GmvcContext) {
		defer func() {
			if x := recover(); x != nil {
				if gmvc.recover != nil {
					ret, err := gmvc.recover(ctx, x)
					if err != nil {
						// TODO: check if has responsed.
						resp := gmvc.errHandler(ctx, err)
						gmvc.doResponse(ctx, resp)
						return
					}

					gmvc.doResponse(ctx, ret)
				} else {
					ctx.Status(http.StatusInternalServerError)
				}
			}
		}()

		ctx.SetActionMeta(actionMeta)
		resp, err := next(ctx)
		if err != nil {
			resp := gmvc.errHandler(ctx, err)
			gmvc.doResponse(ctx, resp)
			return
		}

		gmvc.doResponse(ctx, resp)
	}

	return handlerfunc
}

func (gmvc *GmvcBuilder) Actions() map[string]HandlerFunc {
	return gmvc.actions
}

func (instance *GmvcBuilder) initialize(ctx GmvcContext, handler interface{}) error {
	if initer, ok := handler.(Initer); ok {
		return initer.Init()
	}

	return nil
}

func (instance *GmvcBuilder) launch(ctx GmvcContext, handler interface{}) (interface{}, error) {
	if entity, ok := handler.(Action); ok {
		return entity.Go()
	}

	// this should never happen.
	panic("handler is not implements 'gmvc.Action'.")
}

func (gmvc *GmvcBuilder) doResponse(ctx GmvcContext, resp interface{}) {
	/* if resp == nil, means gmvc no need do response for user */
	if resp == nil {
		return
	}

	if entity, ok := resp.(Response); ok {
		if responsor, ok := gmvc.responsor[entity.Render]; ok {
			responsor.Response(ctx, &entity)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	if entity, ok := resp.(*Response); ok {
		if responsor, ok := gmvc.responsor[entity.Render]; ok {
			responsor.Response(ctx, entity)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}

		return
	}

	// 默认使用JSON进行返回
	gmvc.doResponse(ctx, &Response{
		Render:     JSON,
		Body:       resp,
		StatusCode: http.StatusOK,
	})
}

func (gmvc *GmvcBuilder) resolve(c GmvcContext, meta *ActionMeta) (interface{}, error) {
	// 实例化handler
	handlerValuePtr := reflect.New(meta.handlerType)

	// 解析每个结构体的值
	if err := gmvc.resolveFieldValue(c, handlerValuePtr, meta); err != nil {
		return nil, err
	}

	c.SetAction(handlerValuePtr.Interface() /* the instance pointer of the Action */)

	// check every params
	if err := gmvc.checkFieldValue(c, meta); err != nil {
		return nil, err
	}

	return handlerValuePtr.Interface(), nil
}

func (gmvc *GmvcBuilder) checkFieldValue(ctx GmvcContext, meta *ActionMeta) error {
	for i := 0; i < meta.fieldNum; i++ {
		fieldMeta := meta.fieldList[i]

		if err := gmvc.validateValue(ctx, fieldMeta.value, fieldMeta, meta); err != nil {
			return err
		}
	}

	return nil
}

func (gmvc *GmvcBuilder) resolveFieldValue(ctx GmvcContext, pvalue reflect.Value, meta *ActionMeta) error {
	for i := 0; i < meta.fieldNum; i++ {
		fieldMeta := meta.fieldList[i]

		/* if autowire, try to find singleton instance */
		if len(fieldMeta.autowire) != 0 {
			singleton := gmvc.singletons.GetObjectByName(fieldMeta.autowire)
			if singleton != nil && singleton.typ.AssignableTo(fieldMeta.fieldType.Type) {
				pvalue.Elem().Field(i).Set(reflect.ValueOf(singleton.obj))
			}

			continue
		}

		if fieldMeta.isRecursive {
			recursiveValueStr := reflect.New(fieldMeta.handlerMeta.handlerType)
			if resp := gmvc.resolveFieldValue(ctx, recursiveValueStr, fieldMeta.handlerMeta); resp != nil {
				return resp
			}
			pvalue.Elem().Field(i).Set(recursiveValueStr.Elem())
		}

		var value interface{}

		if reflect.TypeOf(ctx).AssignableTo(fieldMeta.fieldType.Type) {
			value = ctx
		} else if reflect.TypeOf(ctx.GetEntity()).AssignableTo(fieldMeta.fieldType.Type) {
			value = ctx.GetEntity()
		} else {
			originValue, src, ok := gmvc.drawOutOriginValue(ctx, fieldMeta)

			/*
			 found parameter
			  - report param found.
			  - if param found in context, just assign to value.
			  - if param has dedicated resolver, then use the return of resolver as the value.
			  - if param's type match one of the registered typed-resolver, then use typed-resolver to resolve the value.
			  - auto type-convert.
			*/
			if ok {
				ctx.Report(fieldMeta.fieldName)

				if src == CtxSrc {
					value = originValue
				} else if fieldMeta.resolver != nil {
					if tmp, ok := originValue.([]byte); ok {
						originValue = string(tmp)
					}

					var err error
					if value, err = fieldMeta.resolver(ctx, fieldMeta, originValue.(string)); err != nil {
						return err
					}
				} else if resolver, ok := gmvc.typedResolver[fieldMeta.fieldType.Type]; ok {
					if tmp, ok := originValue.([]byte); ok {
						originValue = string(tmp)
					}

					var err error
					if value, err = resolver(ctx, fieldMeta, originValue.(string)); err != nil {
						return err
					}
				} else if src == BodySrc {
					// 处理body, 只有field为string或者[]byte时，才能进行自动赋值
					if fieldMeta.fieldType.Type.Kind() == reflect.String {
						value = string(originValue.([]byte))
					} else {
						value = originValue // must be []byte
					}
				} else {
					value = gmvc.convertFieldValue(ctx, fieldMeta, originValue.(string), meta.handlerName)
				}
			}
		}

		// set resolved value
		fieldMeta.value = value

		if value == nil {
			continue
		}

		valueType := reflect.TypeOf(value)
		if !valueType.AssignableTo(fieldMeta.fieldType.Type) {
			// TODO: warn log
			continue
		}

		// TODO: really need test all type? or just use Set method
		switch valueType.Kind() {
		case reflect.Int:
			pvalue.Elem().Field(i).SetInt(int64(value.(int)))
		case reflect.Int8:
			pvalue.Elem().Field(i).SetInt(int64(value.(int8)))
		case reflect.Int16:
			pvalue.Elem().Field(i).SetInt(int64(value.(int16)))
		case reflect.Int32:
			pvalue.Elem().Field(i).SetInt(int64(value.(int32)))
		case reflect.Int64:
			pvalue.Elem().Field(i).SetInt(value.(int64))
		case reflect.String:
			pvalue.Elem().Field(i).SetString(value.(string))
		case reflect.Uint:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint)))
		case reflect.Uint8:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint8)))
		case reflect.Uint16:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint16)))
		case reflect.Uint32:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint32)))
		case reflect.Uint64:
			pvalue.Elem().Field(i).SetUint(value.(uint64))
		case reflect.Float32:
			pvalue.Elem().Field(i).SetFloat(float64(value.(float32)))
		case reflect.Float64:
			pvalue.Elem().Field(i).SetFloat(value.(float64))
		default:
			pvalue.Elem().Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func (instance *GmvcBuilder) drawOutOriginValue(ctx GmvcContext, fieldMeta *ParamMeta) (originValue interface{}, src Src, present bool) {
	if hasSourceTag(fieldMeta.source, HeaderSrc) {
		src = HeaderSrc
		originValue, present = ctx.GetHeader(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, QuerySrc) {
		src = QuerySrc
		originValue, present = ctx.GetQuery(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, PathSrc) {
		src = PathSrc
		originValue, present = ctx.GetPathParam(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, FormSrc) {
		src = FormSrc
		originValue, present = ctx.GetForm(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, BodySrc) {
		src = BodySrc
		v, err := ctx.GetRawData()
		present = err == nil || len(v) > 0
		originValue = v
		return
	}

	if hasSourceTag(fieldMeta.source, CtxSrc) {
		src = CtxSrc
		originValue, present = ctx.GetCtx(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if fieldMeta.hasDefault {
		src = DefaultSrc
		present = true
		originValue = fieldMeta.def
	}

	return
}

func (instance *GmvcBuilder) validateValue(ctx GmvcContext, value interface{}, fieldMeta *ParamMeta, meta *ActionMeta) error {
	if len(fieldMeta.checkers) <= 0 {
		return nil
	}

	for _, checker := range fieldMeta.checkers {
		if checker == nil {
			continue
		}

		err := checker(ctx, fieldMeta, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// 参数类型转换
func (instance *GmvcBuilder) convertFieldValue(ctx GmvcContext, fieldMeta *ParamMeta, originValue string, handlerName string) interface{} {
	value, err := Convert(originValue, fieldMeta.fieldType.Type)
	if err != nil {
		// logs.Error("Overpass. convert err. handler: %s; value: %s", handlerName, originValue)
		return nil
	}

	return value
}

func (instance *GmvcBuilder) introspect(v reflect.Value) *ActionMeta {
	struct0 := v.Type()

	numField := struct0.NumField()
	structMeta := &ActionMeta{
		handlerName: struct0.Name(),
		handlerType: struct0,
		fieldNum:    numField,
		fieldList:   make([]*ParamMeta, 0),
	}

	for i := 0; i < numField; i++ {
		fieldValue := v.Field(i)
		field := struct0.Field(i)
		tagInfo := field.Tag
		fieldMeta := &ParamMeta{
			fieldType: field,
			tagInfo:   tagInfo,
			fieldName: field.Name,
		}

		// Autowire解析，解析到直接返回，不用继续param的解析
		xAutowire, ok := tagInfo.Lookup(XAutowire)
		if ok {
			fieldMeta.autowire = xAutowire
			fieldMeta.value = fieldValue.Interface()
			structMeta.fieldList = append(structMeta.fieldList, fieldMeta)
			continue
		}

		// xParams解析
		xParams := strings.Split(tagInfo.Get(XParam), XSplit)
		for _, value := range xParams {
			if value == "" {
				continue
			}
			switch value {
			case XRecursive:
				// 如果需要递归解析，则递归下去。
				fieldMeta.isRecursive = true
				fieldMeta.handlerMeta = instance.introspect(fieldValue)
			case XQuery:
				fieldMeta.source |= QuerySrc
			case XForm:
				fieldMeta.source |= FormSrc
			case XBody:
				fieldMeta.source |= BodySrc
			case XHeader:
				fieldMeta.source |= HeaderSrc
				// header 首字母自动大写
				name := fieldMeta.fieldName
				for i, v := range name {
					fieldMeta.fieldName = string(unicode.ToUpper(v)) + name[i+1:]
					break
				}
			case XPath:
				fieldMeta.source |= PathSrc
			case XCtx:
				fieldMeta.source |= CtxSrc
			case XAuto:
				// 如果是'Auto'，则使用Option中的定义
				fieldMeta.source |= Src(instance.options.autodef)
			default:
				fieldMeta.fieldName = value
			}
		}

		// xValidator解析
		validatorStr, ok := tagInfo.Lookup(XValidator)
		if ok {
			xValidator := strings.Split(validatorStr, XSplit)
			checkers := make([]Validator, len(xValidator))
			fieldMeta.checkers = checkers
			for index, value := range xValidator {
				if value == "" {
					continue
				}
				checker := instance.checkerMap[value]
				if checker == nil {
					continue
				}
				checkers[index] = checker
			}
		}

		// Resolver解析
		resolverStr, ok := tagInfo.Lookup(XResolver)
		if ok {
			resolver, ok := instance.resolverMap[resolverStr]
			if ok {
				fieldMeta.resolver = resolver
			}
		}

		// xDefault解析
		defaultStr, ok := tagInfo.Lookup(XDefault)
		if ok {
			fieldMeta.hasDefault = true
			fieldMeta.def = defaultStr
		}

		structMeta.fieldList = append(structMeta.fieldList, fieldMeta)
	}

	return structMeta
}

type singletonContext struct {
	typemap map[reflect.Type]*singleton
	namemap map[string]*singleton
}

func (s *singletonContext) GetObjectByName(name string) *singleton {
	return s.namemap[name]
}

func (s *singletonContext) GetObjectByType(typ reflect.Type) *singleton {
	return s.typemap[typ]
}

type singleton struct {
	name string
	typ  reflect.Type
	obj  interface{}
}
