package v2

import (
	bh "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v2/book"
	"github.com/tchorzewski1991/bds/base/web"
	"go.uber.org/zap"
	"net/http"
)

const version = "v2"

type Config struct {
	Logger *zap.SugaredLogger
}

// Routes binds all the routes for API version 2
func Routes(app *web.App, _ Config) {
	app.Handle(http.MethodGet, version, "/books", bh.List)
}
