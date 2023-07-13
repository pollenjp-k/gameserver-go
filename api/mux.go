package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/clock"
	"github.com/pollenjp/gameserver-go/api/config"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/repository"
	"github.com/pollenjp/gameserver-go/api/service"
)

// multiplexer
func NewMux(ctx context.Context, cfg *config.Config) (
	http.Handler,
	func(), // cleanup func
	error,
) {
	mux := chi.NewRouter()
	mux.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		},
	)

	db, cleanup, err := repository.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	c := clock.RealClocker{}
	r := &repository.Repository{Clocker: c}
	au := auth.NewAuthorizer(db, r)

	{
		v := validator.New()
		cu := &handler.CreateUser{
			Service: &service.CreateUser{
				DB:   db,
				Repo: r,
			},
			Validator: v,
		}
		me := &handler.UserMe{
			Service: &service.GetUser{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		uu := &handler.UpdateUser{
			Service: &service.UpdateUser{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		mux.Route("/user", func(r chi.Router) {
			r.Post("/create", cu.ServeHTTP)
			r.Get("/me", handler.AuthMiddleware(au)(me).ServeHTTP)
			r.Post("/update", handler.AuthMiddleware(au)(uu).ServeHTTP)
		})
	}

	return mux, cleanup, nil
}
