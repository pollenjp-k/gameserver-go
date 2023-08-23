package room

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/service"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out leave_room_moq_test.go . LeaveRoomService
type LeaveRoomService interface {
	LeaveRoom(
		ctx context.Context,
		roomId entity.RoomId,
		userId entity.UserId,
	) error
}

type LeaveRoom struct {
	Service   LeaveRoomService
	Validator *validator.Validate
}

func (wr *LeaveRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		RoomId entity.RoomId `json:"room_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: fmt.Sprintf("decode json: %s", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	if err := wr.Validator.Struct(body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	userId, ok := service.GetUserId(ctx)
	if !ok {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: "failed to get user id from context",
		}, http.StatusInternalServerError)
		return
	}

	if err := wr.Service.LeaveRoom(
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
