package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pollenjp-k/gameserver-go/api/config"
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

	// TODO: not implemented yet
	cleanup := func() {}
	// db, cleanup, err := store.New(ctx, cfg)
	// if err != nil {
	// 	return nil, cleanup, err
	// }
	// c := clock.RealClocker{}
	// r := &store.Repository{Clocker: c}

	// {
	// 	v := validator.New()
	// 	at := &handler.AddTask{
	// 		Service: &service.AddTask{
	// 			DB:   db,
	// 			Repo: r,
	// 		},
	// 		Validator: v,
	// 	}
	// 	lt := &handler.ListTask{
	// 		Service: &service.ListTasks{
	// 			DB:   db,
	// 			Repo: r,
	// 		},
	// 	}
	// 	mux.Route("/tasks", func(r chi.Router) {
	// 		r.Use(handler.AuthMiddleware(jwter))
	// 		r.Post("/", at.ServeHTTP)
	// 		r.Get("/", lt.ServeHTTP)
	// 	})
	// }

	return mux, cleanup, nil
}
