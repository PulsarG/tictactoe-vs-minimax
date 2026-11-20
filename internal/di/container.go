package di

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"tictactoe-vs-minimax/internal/datasource"
	"tictactoe-vs-minimax/internal/domain"
	"tictactoe-vs-minimax/internal/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		providePostgresStorage,
		func(pg *datasource.PostgresStorage) domain.Storage { return pg },
		//fx.Annotate(datasource.NewServerStorage, fx.As(new(domain.Storage))),
		web.NewServer,
		web.NewRouter,
	),
	fx.Invoke(startHTTPServer),
)

func providePostgresStorage(lc fx.Lifecycle) (*datasource.PostgresStorage, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5433/tictactoe?sslmode=disable"
	}
	storage, err := datasource.NewPostgresStorage(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Uber FX: close pool Postgres")
			storage.Close()
			return nil
		},
	})
	return storage, nil
}

func startHTTPServer(lc fx.Lifecycle, r *gin.Engine, srv *web.Server) {
	srv.RegisterRoutes(r)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Uber Fx: server start :8080")
			go func() {
				if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("Error start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Uber Fx: stopped server")
			return httpServer.Shutdown(ctx)
		},
	})
}
