package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/core/book"
	"github.com/tchorzewski1991/bds/business/core/user"
	"github.com/tchorzewski1991/bds/business/web/v1/mid"
	"go.uber.org/zap"
)

const version = "v1"

type Config struct {
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
}

// Routes binds all the routes for API version 1.
func Routes(app *web.App, cfg Config) {
	// Setup book routes.
	bh := bookHandler{book: book.NewCore(cfg.DB, cfg.Logger)}
	app.Handle(http.MethodPost, version, "/books", bh.Create)
	app.Handle(http.MethodGet, version, "/books", bh.Query)
	app.Handle(http.MethodGet, version, "/books/:id", bh.QueryByID)

	// Setup user routes.
	uh := userHandler{user: user.NewCore(cfg.DB, cfg.Logger)}
	app.Handle(http.MethodPost, version, "/user/token", uh.Token)
	app.Handle(http.MethodGet, version, "/user/profile", uh.Profile,
		mid.Authenticate(),
		mid.Authorize("user.profile"),
	)
}
