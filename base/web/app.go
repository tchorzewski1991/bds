package web

import (
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

type App struct {
	mux *httptreemux.ContextMux
}

func NewApp() *App {
	return &App{mux: httptreemux.NewContextMux()}
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
