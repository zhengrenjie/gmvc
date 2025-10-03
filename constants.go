package gmvc

import "math"

type Src int

const (
	// XSplit tag的key，value用英文逗号分隔
	XSplit = ","

	// XParam 指定参数位置、参数属性（可选）、参数名称（可选）
	XParam = "param"

	// XValidator 指定参数绑定的验证器，参数解析成功后，会挨个执行验证器，验证器return错误则直接返回http
	XValidator = "checker"

	// XDefault 参数没传情况下的默认值
	XDefault = "default"

	// XResolver 参数自定义解析插件
	XResolver = "resolver"

	// XAutowire 申明依赖,当前默认singleton
	XAutowire = "autowire"

	// XHeader 从header取参数
	XHeader = "Header"

	// XQuery 从query取参数
	XQuery = "Query"

	// XBody 从body取参数
	XBody = "Body"

	// XForm 从form取参数
	XForm = "Form"

	// XPath 从path取参数
	XPath = "Path"

	// XCtx 从Context取参数，例如gin.Context.Get(key)
	XCtx = "Ctx"

	// XAuto 遍历header、query、body、path、ctx，看是否有匹配的参数，就近原则，以第一个匹配的为主
	XAuto = "Auto"

	// XRecursive 递归解析
	XRecursive = "Recursive"

	// DefaultSrc 参数来源
	DefaultSrc Src = -1

	// HeaderSrc 参数来源
	HeaderSrc Src = 1 << 0

	// QuerySrc 参数来源
	QuerySrc Src = 1 << 1

	// BodySrc 参数来源
	BodySrc Src = 1 << 2

	// PathSrc 参数来源
	PathSrc Src = 1 << 3

	// CtxSrc 参数来源
	CtxSrc Src = 1 << 4

	// 从form中获取
	// multi-form
	// application/x-www-form-urlencoded
	FormSrc Src = 1 << 5

	// Any 参数来源
	AnySrc Src = math.MaxInt32 ^ BodySrc // 异或BodySrc，默认排除从body整体读取
)
