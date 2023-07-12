package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
)

//go:generate go run github.com/matryer/moq -out create_user_moq_test.go . CreateUserService
type CreateUserService interface {
	CreateUser(
		ctx context.Context,
		name string,
		leaderCard entity.LeaderCardIdIDType,
	) (*entity.User, error)
}

type CreateUser struct {
	Service   CreateUserService
	Validator *validator.Validate
}

func (ru *CreateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		Name string `json:"user_name" validate:"required"`
		// 指定されていなければゼロ値として許容する.
		// 今回は値の範囲の明示が無いため.
		LeaderCardId entity.LeaderCardIdIDType `json:"leader_card_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Validator.Struct(body); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
			// Details: []string{
			// 	"failed to validate request body",
			// 	fmt.Sprintf("%+v", body),
			// },
		}, http.StatusBadRequest)
		return
	}

	u, err := ru.Service.CreateUser(ctx, body.Name, body.LeaderCardId)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Token entity.UserTokenType `json:"user_token"`
	}{
		Token: u.Token,
	}
	RespondJson(ctx, w, rsp, http.StatusOK)
}
