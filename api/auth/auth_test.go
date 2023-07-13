package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pollenjp/gameserver-go/api/clock"
	"github.com/pollenjp/gameserver-go/api/entity"
	"github.com/pollenjp/gameserver-go/api/service"
)

func TestAuthorizer(t *testing.T) {
	t.Parallel()

	type want struct {
		errMsg string
	}

	token := entity.UserTokenType(uuid.NewString())
	tests := map[string]struct {
		isOk   bool
		header http.Header
		want   *want
	}{
		"ok": {
			isOk: true,
			header: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
			want: nil,
		},
		"ng_authorization_header": {
			isOk:   false,
			header: http.Header{},
			want: &want{
				errMsg: "authorization header does not exist",
			},
		},
		"ng_bearer_token_format": {
			isOk: false,
			header: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer%s", token)},
			},
			want: &want{
				errMsg: "authorization header format must be 'Bearer <token>'",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		c := clock.FixedClocker{}

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(
				http.MethodGet,
				"/dummy",
				nil,
			)
			key := "Authorization"
			r.Header.Add(key, tt.header.Get(key))

			dummyUser := &entity.User{
				Id:           1,
				Name:         "test",
				Token:        token,
				LeaderCardId: 0,
				CreatedAt:    c.Now(),
				UpdatedAt:    c.Now(),
			}

			moq := &AuthRepositoryMock{}
			moq.GetUserFromTokenFunc = func(
				_ context.Context,
				_ service.Queryer,
				token entity.UserTokenType,
			) (*entity.User, error) {
				if tt.isOk {
					u := dummyUser
					t.Logf("token: %s", token)
					if token != u.Token {
						t.Fatal("token is not match")
					}
					if err := dummyUser.ValidateNotEmpty(); err != nil {
						return nil, err
					}
					return u, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := &Authorizer{
				DB:   nil,
				Repo: moq,
			}
			var err error
			r, err = sut.FillContext(r)
			if err != nil {
				if diff := cmp.Diff(err.Error(), tt.want.errMsg); diff != "" {
					t.Errorf("error is not match (-want +got)\n%s", diff)
				}
				return
			}

			userId, ok := GetUserId(r.Context())
			if !ok {
				// 既に context に user id がセットされているはずなのにセットされていない
				t.Error("user id is not set in context")
			}

			// TODO: この比較は当たり前にパスするから意味がないかも
			if diff := cmp.Diff(userId, dummyUser.Id); diff != "" {
				t.Errorf("user id is not match (-want +got)\n%s", diff)
			}
		})
	}
}
