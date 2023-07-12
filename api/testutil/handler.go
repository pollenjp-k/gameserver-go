package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// 返ってきたデータと期待されるデータのバイト列をJSONに変換して比較する
func AssertJSON(t *testing.T, got, want []byte) {
	t.Helper()

	var gotJSON interface{}
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("got invalid JSON: %v", err)
	}

	var wantJSON interface{}
	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Fatalf("want invalid JSON: %v", err)
	}

	if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("JSON mismatch (-want +got):\n%s", diff)
	}
}

// Response (`got`) が期待される `status` と `body` と同じかどうかを比較する
func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() {
		got.Body.Close()
	})

	gotBody, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if got.StatusCode != status {
		t.Fatalf("status code mismatch: want %d, got %d", status, got.StatusCode)
	}

	if len(gotBody) == 0 && len(body) == 0 {
		// Both response bodies are empty
		return
	}

	AssertJSON(t, gotBody, body)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	t.Cleanup(func() {
		f.Close()
	})

	body, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	return body
}
