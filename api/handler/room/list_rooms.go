package room

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/config"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/service"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out list_room_moq_test.go . GetRoomListService
type GetRoomListService interface {
	GetRoomList(
		ctx context.Context,
		LiveId entity.LiveId,
	) ([]*service.RoomInfoItem, error)
}

type GetRoomList struct {
	Service   GetRoomListService
	Validator *validator.Validate
}

func (ru *GetRoomList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		// 0 の可能性があるので `validate:"required"` はつけない
		LiveId entity.LiveId `json:"live_id"`
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

	rooms, err := ru.Service.GetRoomList(ctx, body.LiveId)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	type item struct {
		RoomId          entity.RoomId `json:"room_id"`
		LiveId          entity.LiveId `json:"live_id"`
		JoinedUserCount int           `json:"joined_user_count"`
		MaxUserCount    int           `json:"max_user_count"`
	}

	roomInfoList := make([]*item, len(rooms))
	for i, roomInfo := range rooms {
		roomInfoList[i] = &item{
			RoomId:          roomInfo.RoomId,
			LiveId:          roomInfo.LiveId,
			JoinedUserCount: roomInfo.JoinedUserCount,
			MaxUserCount:    config.MaxUserCount,
		}
	}

	rsp := struct {
		RoomInfoList []*item `json:"room_info_list"`
	}{
		RoomInfoList: roomInfoList,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
