package gmvc

import (
	"fmt"
	"time"
)

type Next func(ctx GmvcContext) (interface{}, error)

// IMiddleware defines the interface of gmvc middleware.
// It is sugguested to use only Before&After methods or Around method  as Around can fully cover Before&After‘s functionality.
// When you only want to do something before OR after the Action, use Before (After) method.
// When you need to do something before AND after the Action at the same time, use Around method.
// If you used the three methods in the same time, the invoking order will be: Before -> Around -> After,
// which will confuse you somethime when you couldn't figure out what exactlly the excution order it is.
type IMiddleware interface {

	// Before:
	// Invoking before the following middleware (if exist) and action's excution.
	// If any argument of return is not nil, the whole action return.
	Before(ctx GmvcContext) (interface{}, error)

	// After:
	// Invoking after the preceding middleware (if exist) and action's excution.
	// If any argument of return is not nil, the whole action return.
	After(ctx GmvcContext, result interface{}, err error) (interface{}, error)

	// Around:
	// Invoking before and after the preceding middleware (if exist) and action's excution.
	// If any argument of return is not nil, the whole action return.
	Around(ctx GmvcContext, next Next) (interface{}, error)

	// IsApply determines whether the middleware should be applied to the current request.
	// If it returns false, the middleware will be skipped.
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
