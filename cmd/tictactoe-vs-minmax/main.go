package main

import (
	"net/http"
	"tictactoe-vs-minmax/internal/di"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		di.Module,
		fx.Invoke(func(mux *http.ServeMux) {
			go http.ListenAndServe(":8080", mux)
		}),
	)
	app.Run()
}
