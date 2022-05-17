package handler

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	perrors "github.com/pkg/errors"
)

/**
修改自 Kratos.recovery, 搭配 pkg/errors 和 自定义 zero_kratos 转接器实现 kratos 错误的格式化输出
*/

// ErrUnknownRequest is unknown request error.
var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
	logger  log.Logger
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger with recovery logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) middleware.Middleware {
	op := options{
		logger: log.GetLogger(),
		handler: func(ctx context.Context, req, err interface{}) error {
			return ErrUnknownRequest
		},
	}
	for _, o := range opts {
		o(&op)
	}
	logger := op.logger
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					_ = logger.Log(log.LevelError, "err", perrors.WithStack(rerr.(error)), "req", req)
					err = op.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
