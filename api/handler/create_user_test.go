package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/testutil"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/create_user/ok/req.json.golden",
			want: want{
				status:  200, // http.StatusOK
				rspFile: "testdata/create_user/ok/res.json.golden",
			},
		},
		"bad_empty_username": {
			reqFile: "testdata/create_user/bad_empty_username/req.json.golden",
			want: want{
				status:  400, // http.StatusBadRequest
				rspFile: "testdata/create_user/bad_empty_username/res.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/user/create",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			// auto generated moq
			moq := &CreateUserServiceMock{}
			moq.CreateUserFunc = func(
				ctx context.Context,
				name string,
				leaderCard entity.LeaderCardIdIDType,
			) (*entity.User, error) {
				if tt.want.status == http.StatusOK {
					return &entity.User{
						Id:    1,
						Token: "7e13ae4f-bd7f-4c2a-ad90-8c2c7bf6f425",
					}, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := CreateUser{
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
