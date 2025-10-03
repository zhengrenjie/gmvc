package gmvc

import (
	"fmt"
	"time"
)

type Next func(ctx GmvcContext) (interface{}, error)

// IMiddleware defines the interface of gmvc middleware.
// Only using one of Before&After methods or Around method is sugguested as Around can fully cover Before&After‘s functionality.
// When you only want to do something before OR after Action, use Before (After) method.
// When you need to do something before AND after Action at the same time, use Around method.
// If you used the three methods in the same time, the invoking order will be: Before -> Around -> After,
// which will confuse you somethime when you couldn't figure out what exactlly the excution order it is.
type IMiddleware interface {
	// Before
	// Invoking before the following middleware (if exist) and action's excution.
	// If any argument of return is not nil, the whole action return.
	Before(ctx GmvcContext) (interface{}, error)

	// After
	// 在一个Action执行后执行
	// 入参是前面执行得到的result和err
	// 返回是更新后（或者没有更新）的result和err
	After(ctx GmvcContext, result interface{}, err error) (interface{}, error)

	// Around
	// 在一个Action前后分别执行一些动作
	// 调用next(ctx)执行后续的步骤，通过next(ctx)获得后续执行的结果
	// Next会执行后续所有的Middleware和Action需要执行的步骤，并返回前面所有执行后的返回结果
	Around(ctx GmvcContext, next Next) (interface{}, error)

	// IsApply 动态判断是否执行当前middleware
	IsApply(ctx GmvcContext) bool
}

var _ IMiddleware = (*BaseMiddleware)(nil)

// BaseMiddleware defines the basic action of a Middleware.
// It has empty Before & Around & After methods.
// Embedded this strcut in your own Middleware is suggested as you only need to implement the methods you actually need.
type BaseMiddleware struct {
}

// IsApply implements IMiddleware.
func (b *BaseMiddleware) IsApply(ctx GmvcContext) bool {
	return true // 默认开启本middleware
}

// After implements IMiddleware.
// Empty implements.
func (*BaseMiddleware) After(ctx GmvcContext, result interface{}, err error) (interface{}, error) {
	return result, err
}

// Around implements IMiddleware.
// Empty implements.
func (*BaseMiddleware) Around(ctx GmvcContext, next Next) (interface{}, error) {
	return next(ctx)
}

// Before implements IMiddleware.
// Empty implements.
func (*BaseMiddleware) Before(ctx GmvcContext) (interface{}, error) {
	return nil, nil
}

// TODO: delete before merge
type ExampleMiddleware struct {
	BaseMiddleware
}

func (*ExampleMiddleware) Around(ctx GmvcContext, next Next) (interface{}, error) {
	begin := time.Now()

	ret, err := next(ctx)

	end := time.Now()

	fmt.Printf("duration %d", end.Sub(begin))

	return ret, err
}
