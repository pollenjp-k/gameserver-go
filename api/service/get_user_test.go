package service

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pollenjp/gameserver-go/api/clock"
	"github.com/pollenjp/gameserver-go/api/entity"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		isErr  bool
		errMsg string
	}{
		"ok": {
			isErr:  false,
			errMsg: "",
		},
		"ng_user_not_found": {
			isErr:  true,
			errMsg: "",
		},
	}

	for n, tt := range tests {
		tt := tt
		c := clock.FixedClocker{}
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			dummyUser := &entity.User{
				Id:           1,
				Name:         "test",
				Token:        entity.UserTokenType(uuid.NewString()),
				LeaderCardId: 1,
				CreatedAt:    c.Now(),
				UpdatedAt:    c.Now(),
			}

			moq := &UserGetterMock{}
			moq.GetUserFromIdFunc = func(
				_ context.Context,
				_ Queryer,
				_ entity.UserId,
			) (*entity.User, error) {
				u := dummyUser
				if err := u.ValidateNotEmpty(); err != nil {
					return nil, err
				}
				return u, nil
			}

			sut := &GetUser{
				DB:   nil,
				Repo: moq,
			}
			u, err := sut.GetUser(
				context.Background(),
				dummyUser.Id,
			)
			if err != nil {
				if !tt.isErr {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// エラーメッセージの一致を期待
				if diff := cmp.Diff(tt.errMsg, err.Error()); diff != "" {
					t.Errorf("error message mismatch (-want +got):\n%s", diff)
				}
				return
			}

			// ユーザー情報の一致を期待
			if diff := cmp.Diff(dummyUser, u); diff != "" {
				t.Errorf("user mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
