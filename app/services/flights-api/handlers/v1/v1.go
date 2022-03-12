package v1

import (
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the version 1 routes
func Routes(mux http.Handler, cfg Config) {}
