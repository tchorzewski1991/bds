// Package logger provides a helper function responsible for encapsulating logger construction.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(fields ...Field) (*zap.SugaredLogger, error) {
	conf := zap.NewProductionConfig()
	conf.DisableStacktrace = true
	conf.DisableCaller = true
	conf.OutputPaths = []string{"stdout"}
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Assign initial fields common to all of logs
	initialFields := make(map[string]interface{})
	for _, f := range fields {
		initialFields[f.Name] = f.Value
	}
	conf.InitialFields = initialFields

	// Build logger from config
	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

type Field struct {
	Name  string
	Value interface{}
}
