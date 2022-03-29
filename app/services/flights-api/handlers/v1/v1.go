package v1

import (
	fh "github.com/tchorzewski1991/fds/app/services/flights-api/handlers/v1/flight"
	uh "github.com/tchorzewski1991/fds/app/services/flights-api/handlers/v1/user"
	"github.com/tchorzewski1991/fds/base/web"
	"github.com/tchorzewski1991/fds/business/web/v1/mid"
	"go.uber.org/zap"
	"net/http"
)

const version = "v1"

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the routes for API version 1
func Routes(app *web.App, cfg Config) {
	// Flight handlers
	app.Handle(http.MethodGet, version, "/flights", fh.List)
	app.Handle(http.MethodGet, version, "/flights/:id", fh.QueryByID)

	// User handlers
	app.Handle(http.MethodPost, version, "/user/token", uh.Token)
	app.Handle(http.MethodGet, version, "/user/protected", uh.Protected, mid.Authenticate())
}
