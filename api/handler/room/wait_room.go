package room

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/auth"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/service"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out wait_room_moq_test.go . WaitRoomService
type WaitRoomService interface {
	WaitRoom(
		ctx context.Context,
		roomId entity.RoomId,
		userId entity.UserId,
	) (*service.WaitRoomResult, error)
}

type WaitRoom struct {
	Service   WaitRoomService
	Validator *validator.Validate
}

func (wr *WaitRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err := wr.Validator.Struct(body); err != nil {
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

	waitRoomResult, err := wr.Service.WaitRoom(
		ctx,
		body.RoomId,
		userId,
	)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		Status       entity.RoomStatus          `json:"status"`
		RoomInfoList []*service.WaitingRoomUser `json:"room_user_list"`
	}{
		Status:       waitRoomResult.Room.Status,
		RoomInfoList: waitRoomResult.WaitingRoomUser,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
