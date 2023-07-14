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
	"github.com/pollenjp/gameserver-go/api/handler/room"
	"github.com/pollenjp/gameserver-go/api/handler/user"
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
		cu := &user.CreateUser{
			Service: &service.CreateUser{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		me := &user.UserMe{
			Service: &service.GetUser{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		uu := &user.UpdateUser{
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

	{
		cr := &room.CreateRoom{
			Service: &service.CreateRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		rl := &room.GetRoomList{
			Service: &service.GetRoomList{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		jr := &room.JoinRoom{
			Service: &service.JoinRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		wr := &room.WaitRoom{
			Service: &service.WaitRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		sr := &room.StartRoom{
			Service: &service.StartRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		er := &room.EndRoom{
			Service: &service.EndRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		rr := &room.RoomResult{
			Service: &service.GetRoomResult{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		lr := &room.LeaveRoom{
			Service: &service.LeaveRoom{
				DB:   db,
				Repo: r,
			},
			Validator: validator.New(),
		}
		mux.Route("/room", func(r chi.Router) {
			r.Post("/create", handler.AuthMiddleware(au)(cr).ServeHTTP)
			r.Post("/list", rl.ServeHTTP)
			r.Post("/join", handler.AuthMiddleware(au)(jr).ServeHTTP)
			r.Post("/wait", handler.AuthMiddleware(au)(wr).ServeHTTP)
			r.Post("/start", handler.AuthMiddleware(au)(sr).ServeHTTP)
			r.Post("/end", handler.AuthMiddleware(au)(er).ServeHTTP)
			r.Post("/result", handler.AuthMiddleware(au)(rr).ServeHTTP)
			r.Post("/leave", handler.AuthMiddleware(au)(lr).ServeHTTP)
		})
	}

	return mux, cleanup, nil
}
