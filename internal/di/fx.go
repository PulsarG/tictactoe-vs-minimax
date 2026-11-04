package di

import (
	"net/http"
	"tictactoe-vs-minmax/internal"
	"tictactoe-vs-minmax/internal/app"
	"tictactoe-vs-minmax/internal/web"

	"go.uber.org/fx"
)

var Module = fx.Option(
	fx.Provide(
		internal.NewStorage,
		internal.NewGameRepository,
		app.NewGameService,
		web.NewGameHandler,
	),
	fx.Invoke(registerHTTP),
)

func registerHTTP(handler *web.GameHandler) *http.ServeMux {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}
