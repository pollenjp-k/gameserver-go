package room

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/service"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out join_room_moq_test.go . JoinRoomService
type JoinRoomService interface {
	JoinRoom(
		ctx context.Context,
		roomId entity.RoomId,
		userId entity.UserId,
		liveDifficulty entity.LiveDifficulty,
	) (entity.JoinRoomResult, error)
}

type JoinRoom struct {
	Service   JoinRoomService
	Validator *validator.Validate
}

func (ru *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		RoomId           entity.RoomId         `json:"room_id" validate:"required"`
		SelectDifficulty entity.LiveDifficulty `json:"select_difficulty" validate:"required"`
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

	userId, ok := service.GetUserId(ctx)
	if !ok {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: "failed to get user id from context",
		}, http.StatusInternalServerError)
		return
	}

	result, err := ru.Service.JoinRoom(
		ctx,
		body.RoomId,
		userId,
		body.SelectDifficulty,
	)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		JoinRoomResult entity.JoinRoomResult `json:"join_room_result"`
	}{
		JoinRoomResult: result,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
