package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pollenjp/gameserver-go/api/config"
	"github.com/pollenjp/gameserver-go/api/entity"
	userHandler "github.com/pollenjp/gameserver-go/api/handler/user"
)

// response body from `/user/create`
func GotBodyOfUserCreate(t *testing.T, mux http.Handler, reqBody userHandler.CreateUserRequestJson) []byte {
	t.Helper()

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/create", bytes.NewBuffer(reqBodyJson))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	rsp := w.Result()
	defer func() {
		_ = rsp.Body.Close()
	}()

	gotBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
	}

	return gotBody
}

// - `/user/create`
func TestNewMuxUserCreate(t *testing.T) {
	t.Parallel()

	// setup
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
		Name:         "test",
		LeaderCardId: 1,
	})

	// ありのままに変換
	{
		var gotJson interface{}
		if err := json.Unmarshal([]byte(gotBody), &gotJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		// 期待するキーの数を確認
		expectedKeyNum := 1
		if len(gotJson.(map[string]interface{})) != expectedKeyNum {
			t.Errorf("expected to have %d key, but got %d", expectedKeyNum, len(gotJson.(map[string]interface{})))
		}
	}

	// 期待する型に変換
	{
		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
	}
}

// - `/user/me`
func TestNewMuxUserMe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	// `/user/create`
	{
		gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		})

		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		sampleUser.Token = gotTypedJson.Token
	}

	// `/user/me`
	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sampleUser.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		if rsp.StatusCode != http.StatusOK {
			t.Fatalf("status code (want %d, got %d)", http.StatusOK, rsp.StatusCode)
		}

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		var gotTypedJson userHandler.UserMeResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		expected := userHandler.UserMeResponseJson{
			Id:           gotTypedJson.Id,
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		}

		if diff := cmp.Diff(expected, gotTypedJson); diff != "" {
			t.Fatalf("diff: (-want +got)\n%s", diff)
		}
	}
}

// - `/user/update`
func TestNewMuxUserUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	// `/user/create`
	{
		gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		})

		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		sampleUser.Token = gotTypedJson.Token
	}

	updatedUser := entity.User{
		Name:         "updated",
		LeaderCardId: 2,
		Token:        sampleUser.Token,
	}

	// `/user/update`
	{
		reqJsonBody := []byte(
			fmt.Sprintf(
				`{"user_name":"%s","leader_card_id":%d}`,
				updatedUser.Name,
				updatedUser.LeaderCardId,
			),
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/update", bytes.NewBuffer(reqJsonBody))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", updatedUser.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if rsp.StatusCode != http.StatusOK {
			FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
		}

		var gotJson interface{}
		if err := json.Unmarshal([]byte(gotBody), &gotJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		expectedKeyNum := 0
		if len(gotJson.(map[string]interface{})) != expectedKeyNum {
			t.Errorf("expected to have %d key, but got %d", expectedKeyNum, len(gotJson.(map[string]interface{})))
		}
	}
}

func FatalErrorWithStatusCodeAndBody(t *testing.T, expectedStatusCode int, gotStatusCode int, gotBody []byte) {
	t.Helper()

	t.Errorf("status code (want %d, got %d)", expectedStatusCode, gotStatusCode)
	var errorJson interface{}
	if err := json.Unmarshal([]byte(gotBody), &errorJson); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	t.Fatalf("error json:%v", errorJson)
}
