package v2

import (
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the version 2 routes
func Routes(mux http.Handler, cfg Config) {}
