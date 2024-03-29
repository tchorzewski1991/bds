package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/sys/metrics"
	"go.uber.org/zap"
)

type ExtContext struct {
	extContext   sqlx.ExtContext
	interceptors []interceptor
}

type Opt func(extContext *ExtContext)

type ErrorMapper func(err error) error

type Metric interface {
	Send()
}

func NewExtContext(extContext sqlx.ExtContext) *ExtContext {
	return &ExtContext{extContext: extContext}
}

func NewMetric(table, operation string) Metric {
	return metrics.DbHistogram(table, operation)
}

func NewErrorMapper() ErrorMapper {
	return func(err error) error {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
}

func (ec *ExtContext) WithLogger(logger *zap.SugaredLogger) *ExtContext {
	ec.interceptors = append(ec.interceptors, loggerInterceptor(logger))
	return ec
}

func (ec *ExtContext) WithMetric(metric Metric) *ExtContext {
	ec.interceptors = append(ec.interceptors, metricInterceptor(metric))
	return ec
}

func (ec *ExtContext) WithErrorMapper(mapper ErrorMapper) *ExtContext {
	ec.interceptors = append(ec.interceptors, errorMapperInterceptor(mapper))
	return ec
}

func loggerInterceptor(logger *zap.SugaredLogger) interceptor {

	i := func(handler handler) handler {

		h := func(ctx context.Context, query string, args ...any) (any, error) {
			logger.Infow("db call", "trace_id", web.GetTraceID(ctx), "query", query, "args", args)
			return handler(ctx, query, args...)
		}

		return h

	}

	return i
}

func metricInterceptor(metric Metric) interceptor {

	i := func(handler handler) handler {

		h := func(ctx context.Context, query string, args ...any) (any, error) {
			defer metric.Send()
			return handler(ctx, query, args...)
		}

		return h

	}

	return i
}

func errorMapperInterceptor(mapper ErrorMapper) interceptor {

	i := func(handler handler) handler {

		h := func(ctx context.Context, query string, args ...any) (any, error) {

			result, err := handler(ctx, query, args...)
			if err != nil {
				err = mapper(err)
				return nil, err
			}

			return result, nil
		}

		return h

	}

	return i
}

func (ec *ExtContext) DriverName() string {
	return ec.extContext.DriverName()
}

func (ec *ExtContext) Rebind(query string) string {
	return ec.extContext.Rebind(query)
}

func (ec *ExtContext) BindNamed(query string, args any) (string, []any, error) {
	return ec.extContext.BindNamed(query, args)
}

type interceptor func(handler) handler

type handler func(ctx context.Context, query string, args ...any) (any, error)

func wrapInterceptors(interceptors []interceptor, handler handler) handler {
	for i := len(interceptors) - 1; i >= 0; i-- {
		intc := interceptors[i]
		if intc != nil {
			handler = intc(handler)
		}
	}
	return handler
}

func (ec *ExtContext) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	h := wrapInterceptors(ec.interceptors, func(ctx context.Context, query string, args ...any) (any, error) {
		return ec.extContext.QueryContext(ctx, query, args...)
	})
	result, err := h(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return result.(*sql.Rows), nil
}

func (ec *ExtContext) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	h := wrapInterceptors(ec.interceptors, func(ctx context.Context, query string, args ...any) (any, error) {
		return ec.extContext.QueryxContext(ctx, query, args...)
	})
	result, err := h(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return result.(*sqlx.Rows), nil
}

func (ec *ExtContext) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return ec.extContext.QueryRowxContext(ctx, query, args...)
}

func (ec *ExtContext) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	h := wrapInterceptors(ec.interceptors, func(ctx context.Context, query string, args ...any) (any, error) {
		return ec.extContext.ExecContext(ctx, query, args...)
	})
	result, err := h(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return result.(sql.Result), nil
}
