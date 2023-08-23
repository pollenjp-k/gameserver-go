package room

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/handler"
	"github.com/pollenjp/gameserver-go/api/service"
)

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out end_room_moq_test.go . EndRoomService
type EndRoomService interface {
	EndRoom(
		ctx context.Context,
		score *entity.Score,
	) error
}

type EndRoom struct {
	Service   EndRoomService
	Validator *validator.Validate
}

func (ru *EndRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body struct {
		RoomId         entity.RoomId `json:"room_id" validate:"required"`
		Score          int           `json:"score" validate:"required"`
		JudgeCountList []int         `json:"judge_count_list" validate:"required,list_length=5"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: fmt.Sprintf("decode json: %s", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Validator.RegisterValidation("list_length", func(fl validator.FieldLevel) bool {
		param := fl.Param()
		expectedLength, err := strconv.Atoi(param)
		if err != nil {
			return false
		}

		// Field().Interface()でinterface{}型を取得し、[]intにキャスト
		value, _ := fl.Field().Interface().([]int)
		return len(value) == expectedLength
	}); err != nil {
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

	score := entity.NewScore(
		body.RoomId,
		userId,
		body.Score,
		body.JudgeCountList[0],
		body.JudgeCountList[1],
		body.JudgeCountList[2],
		body.JudgeCountList[3],
		body.JudgeCountList[4],
	)

	if err := ru.Service.EndRoom(
		ctx,
		score,
	); err != nil {
		handler.RespondJson(ctx, w, &handler.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct{}{}
	handler.RespondJson(ctx, w, rsp, http.StatusOK)
}
