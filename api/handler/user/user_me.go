package user

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
)

//go:generate go run github.com/matryer/moq -out usre_me_moq_test.go . GetUserService
type GetUserService interface {
	GetUser(
		ctx context.Context,
		userId entity.UserId,
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
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: "failed to get user id from token",
		}, http.StatusInternalServerError)
		return
	}

	u, err := ru.Service.GetUser(ctx, userId)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Id           entity.UserId             `json:"id"`
		Name         string                    `json:"name"`
		LeaderCardId entity.LeaderCardIdIDType `json:"leader_card_id"`
	}{
		Id:           u.Id,
		Name:         u.Name,
		LeaderCardId: u.LeaderCardId,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
