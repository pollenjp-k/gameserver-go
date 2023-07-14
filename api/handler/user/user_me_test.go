package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pollenjp/gameserver-go/api/clock"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
	"github.com/pollenjp/gameserver-go/api/testutil"
)

func TestUserMe(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		want want
	}{
		"ok": {
			want: want{
				status:  200, // http.StatusOK
				rspFile: "testdata/user_me/ok/res.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		c := clock.FixedClocker{}
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/user/me",
				nil,
			)

			dummyUser := &entity.User{
				Id:           1,
				Name:         "test",
				Token:        entity.UserTokenType(uuid.NewString()),
				LeaderCardId: 0,
				CreatedAt:    c.Now(),
				UpdatedAt:    c.Now(),
			}

			moq := &GetUserServiceMock{}
			moq.GetUserFunc = func(
				_ context.Context,
				_ entity.UserId,
			) (*entity.User, error) {
				if tt.want.status == http.StatusOK {
					u := dummyUser
					if err := u.ValidateNotEmpty(); err != nil {
						return nil, err
					}
					return u, nil
				}
				return nil, errors.New("error from mock")
			}

			// 認証情報の追加
			ctx := service.SetUserId(r.Context(), dummyUser.Id)
			r = r.Clone(ctx)

			sut := UserMe{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(
				t,
				rsp,
				tt.want.status,
				testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
