package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
)

//go:generate go run github.com/matryer/moq -out usre_me_moq_test.go . UpdateUserService
type UpdateUserService interface {
	UpdateUser(
		ctx context.Context,
		user *entity.User,
	) error
}

type UpdateUser struct {
	Service   UpdateUserService
	Validator *validator.Validate
}

func (ru *UpdateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, isOk := auth.GetUserId(ctx)
	if !isOk {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: "failed to get user id from token",
		}, http.StatusInternalServerError)
		return
	}

	var body struct {
		Name         string                    `json:"user_name" validate:"required"`
		LeaderCardId entity.LeaderCardIdIDType `json:"leader_card_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Validator.Struct(body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	err := ru.Service.UpdateUser(ctx, &entity.User{
		Id:           userId,
		Name:         body.Name,
		LeaderCardId: body.LeaderCardId,
	})
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct{}{}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
