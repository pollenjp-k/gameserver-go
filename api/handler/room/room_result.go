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
// go:generate go run github.com/matryer/moq -out room_result_moq_test.go . RoomResultService
type RoomResultService interface {
	GetRoomResult(
		ctx context.Context,
		roomId entity.RoomId,
	) ([]*service.RoomUserResult, error)
}

type RoomResult struct {
	Service   RoomResultService
	Validator *validator.Validate
}

func (ru *RoomResult) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			// Details: []string{
			// 	"failed to validate request body",
			// 	fmt.Sprintf("%+v", body),
			// },
		}, http.StatusBadRequest)
		return
	}

	roomUserResults, err := ru.Service.GetRoomResult(ctx, body.RoomId)
	if err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	type item struct {
		UserId         entity.UserId `json:"user_id"`
		Score          int           `json:"score"`
		JudgeCountList []int         `json:"judge_count_list"`
	}

	ResultUserList := make([]*item, len(roomUserResults))
	for i, roomInfo := range roomUserResults {
		ResultUserList[i] = &item{
			UserId: roomInfo.UserId,
			Score:  roomInfo.Score,
			JudgeCountList: []int{
				roomInfo.JudgePerfect,
				roomInfo.JudgeGreat,
				roomInfo.JudgeGood,
				roomInfo.JudgeBad,
				roomInfo.JudgeMiss,
			},
		}
	}

	rsp := struct {
		ResultUserList []*item `json:"result_user_list"`
	}{
		ResultUserList: ResultUserList,
	}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
