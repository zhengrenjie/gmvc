package gmvc

import "context"

var logger ILogger

func SetLogger(log ILogger) {
	logger = log
}

type ILogger interface {
	Debug(ctx context.Context, msg string, vars ...interface{})
	Info(ctx context.Context, msg string, vars ...interface{})
	Error(ctx context.Context, msg string, vars ...interface{})
}
