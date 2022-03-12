package v2

import (
	"github.com/tchorzewski1991/fds/base/web"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the version 2 routes
func Routes(app *web.App, cfg Config) {}
