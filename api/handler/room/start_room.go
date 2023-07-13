package room

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out start_room_moq_test.go . StartRoomService
type StartRoomService interface {
	StartRoom(
		ctx context.Context,
		roomId entity.RoomId,
		userId entity.UserId,
	) error
}

type StartRoom struct {
	Service   StartRoomService
	Validator *validator.Validate
}

func (ru *StartRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		RoomId entity.RoomId `json:"room_id" validate:"required"`
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

	userId, ok := auth.GetUserId(ctx)
	if !ok {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: "failed to get user id from context",
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Service.StartRoom(
		ctx,
		body.RoomId,
		userId,
	); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct{}{}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
