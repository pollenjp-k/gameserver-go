package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
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

type CreateUserRequestJson struct {
	Name         string                    `json:"user_name" validate:"required"`
	LeaderCardId entity.LeaderCardIdIDType `json:"leader_card_id" validate:"required"`
}

type CreateUserResponseJson struct {
	Token entity.UserTokenType `json:"user_token"`
}

func (ru *CreateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body CreateUserRequestJson
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: fmt.Sprintf("decode json: %s", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Validator.Struct(body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
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
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := CreateUserResponseJson{
		Token: u.Token,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
