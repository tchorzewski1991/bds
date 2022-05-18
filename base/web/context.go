package web

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type ctxKey int

const key ctxKey = 1

type CtxValues struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func GetCtxValues(ctx context.Context) (*CtxValues, error) {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return nil, errors.New("ctx values not present")
	}
	return v, nil
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return errors.New("ctx values not present")
	}
	v.StatusCode = statusCode
	return nil
}

func GetStatusCode(ctx context.Context) string {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return "000"
	}
	return strconv.Itoa(v.StatusCode)
}
