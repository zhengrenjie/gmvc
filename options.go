package gmvc

type GmvcOptions struct {
	autodef Src
}

type GmvcOption func(options GmvcOptions)

// DefineAuto
// 定义Auto的行为，从HTTP协议的哪些地方自动获取参数，例如，Auto=Query|Body，则会自动从Query和Body的地方来获取参数
func DefineAuto(srclist ...Src) GmvcOption {
	return func(options GmvcOptions) {
		var auto Src = 0
		for _, src := range srclist {
			auto |= src
		}

		options.autodef = auto
	}
}
