package handler

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/entity"
)

//go:generate go run github.com/matryer/moq -out usre_me_moq_test.go . GetUserService
type GetUserService interface {
	GetUser(
		ctx context.Context,
		userId entity.UserID,
	) (*entity.User, error)
}

type UserMe struct {
	Service   GetUserService
	Validator *validator.Validate
}

func (ru *UserMe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, isOk := auth.GetUserId(ctx)
	if !isOk {
		RespondJson(ctx, w, &ErrResponse{
			Message: "failed to get user id from token",
		}, http.StatusInternalServerError)
		return
	}

	u, err := ru.Service.GetUser(ctx, userId)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Id           entity.UserID             `json:"id"`
		Name         string                    `json:"name"`
		LeaderCardId entity.LeaderCardIdIDType `json:"leader_card_id"`
	}{
		Id:           u.Id,
		Name:         u.Name,
		LeaderCardId: u.LeaderCardId,
	}
	RespondJson(ctx, w, rsp, http.StatusOK)
}
