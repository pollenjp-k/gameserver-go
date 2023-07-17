package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pollenjp/gameserver-go/api/config"
	userHandler "github.com/pollenjp/gameserver-go/api/handler/user"
)

// response body from `/user/create`
func GotBodyOfUserCreate(t *testing.T, mux http.Handler) []byte {
	t.Helper()

	reqBody := userHandler.CreateUserRequestJson{
		Name:         "test",
		LeaderCardId: 1,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/create", bytes.NewBuffer(reqBodyJson))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	rsp := w.Result()
	t.Cleanup(func() {
		_ = rsp.Body.Close()
	})

	if rsp.StatusCode != http.StatusOK {
		t.Fatalf("status code (want %d, got %d)\n%v", http.StatusOK, rsp.StatusCode, rsp.Body)
	}

	gotBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	return gotBody
}

// - `/user/create`
func TestNewMuxPathUserCreate(t *testing.T) {
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

	gotBody := GotBodyOfUserCreate(t, mux)

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
