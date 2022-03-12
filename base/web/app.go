package web

import (
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"os"
	"syscall"
)

type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		mux:      httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Shutdown() {
	a.shutdown <- syscall.SIGTERM
}
