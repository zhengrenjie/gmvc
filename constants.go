package gmvc

import "math"

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

	// XHeader 从header取参数
	XHeader = "Header"

	// XQuery 从query取参数
	XQuery = "Query"

	// XBody 从body取参数
	XBody = "Body"

	// XPath 从path取参数
	XPath = "Path"

	// XCtx 从Context取参数，例如gin.Context.Get(key)
	XCtx = "Ctx"

	// XAuto 遍历header、query、body、path、ctx，看是否有匹配的参数，就近原则，以第一个匹配的为主
	XAuto = "Auto"

	// XRecursive 递归解析
	XRecursive = "Recursive"

	// DefaultSrc 参数来源
	DefaultSrc = -1

	// HeaderSrc 参数来源
	HeaderSrc = 1 << 0

	// QuerySrc 参数来源
	QuerySrc = 1 << 1

	// BodySrc 参数来源
	BodySrc = 1 << 2

	// PathSrc 参数来源
	PathSrc = 1 << 3

	// CtxSrc 参数来源
	CtxSrc = 1 << 4

	// AutoSrc 参数来源
	AutoSrc = math.MaxInt32
)
