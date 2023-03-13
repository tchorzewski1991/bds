package v2

import (
	"net/http"

	"github.com/tchorzewski1991/bds/base/web"
	"go.uber.org/zap"
)

const version = "v2"

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the routes for API version 2.
func Routes(app *web.App, _ Config) {
	app.Handle(http.MethodGet, version, "/books", List)
}
