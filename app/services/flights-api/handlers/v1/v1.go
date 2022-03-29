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
	app.Handle(http.MethodGet, version, "/flights", fh.List)
	app.Handle(http.MethodGet, version, "/flights/:id", fh.QueryByID)

	// The following endpoint exists just for testing.
	// It will be removed after properly developed authorization mechanism.
	app.Handle(http.MethodGet, version, "/protected", fh.Protected, mid.Authenticate())

	app.Handle(http.MethodPost, version, "/user/token", uh.Token)
}
