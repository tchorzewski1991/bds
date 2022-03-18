package v1

import (
	fh "github.com/tchorzewski1991/fds/app/services/flights-api/handlers/v1/flight"
	"github.com/tchorzewski1991/fds/base/web"
	"go.uber.org/zap"
	"net/http"
)

const version = "v1"

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the routes for API version 1
func Routes(app *web.App, _ Config) {
	app.Handle(http.MethodGet, version, "/flights", fh.List)
	app.Handle(http.MethodGet, version, "/flights/:id", fh.QueryByID)
}
