package web

import (
	"context"
	"errors"
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

func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*CtxValues)
	if !ok {
		return errors.New("ctx values not present")
	}
	v.StatusCode = statusCode
	return nil
}
