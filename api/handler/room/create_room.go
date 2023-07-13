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

//go:generate go run github.com/matryer/moq -out create_room_moq_test.go . CreateRoomService
type CreateRoomService interface {
	CreateRoom(
		ctx context.Context,
		liveId entity.LiveId,
		hostUserId entity.UserId,
	) (*entity.Room, *entity.RoomUser, error)
}

type CreateRoom struct {
	Service   CreateRoomService
	Validator *validator.Validate
}

func (ru *CreateRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		// 0 の可能性がある
		LiveId           entity.LiveId         `json:"live_id"`
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
			// Details: []string{
			// 	"failed to validate request body",
			// 	fmt.Sprintf("%+v", body),
			// },
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

	room, _, err := ru.Service.CreateRoom(ctx, body.LiveId, userId)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		RoomId entity.RoomId `json:"room_id"`
	}{
		RoomId: room.Id,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
