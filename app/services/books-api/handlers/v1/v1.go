package v1

import (
	"github.com/jmoiron/sqlx"
	bh "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v1/book"
	uh "github.com/tchorzewski1991/bds/app/services/books-api/handlers/v1/user"
	"github.com/tchorzewski1991/bds/base/web"
	"github.com/tchorzewski1991/bds/business/core/book"
	"github.com/tchorzewski1991/bds/business/core/user"
	"github.com/tchorzewski1991/bds/business/web/v1/mid"
	"go.uber.org/zap"
	"net/http"
)

const version = "v1"

type Config struct {
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
}

// Routes binds all the routes for API version 1
func Routes(app *web.App, cfg Config) {
	// Book handlers group
	bhg := bh.Handler{Book: book.NewCore(cfg.DB, cfg.Logger)}
	app.Handle(http.MethodPost, version, "/books", bhg.Create)
	app.Handle(http.MethodGet, version, "/books", bhg.Query)
	app.Handle(http.MethodGet, version, "/books/:id", bhg.QueryByID)

	// User handlers group
	uhg := uh.Handler{User: user.NewCore(cfg.DB, cfg.Logger)}
	app.Handle(http.MethodPost, version, "/user/token", uhg.Token)
	app.Handle(http.MethodGet, version, "/user/profile", uhg.Profile,
		mid.Authenticate(),
		mid.Authorize("user.profile"),
	)
}
