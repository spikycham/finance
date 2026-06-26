package network

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadBody(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		body    string
		want    *TestData
		wantErr error
	}{
		{
			name:    "valid JSON body",
			body:    `{"name":"test","value":42}`,
			want:    &TestData{Name: "test", Value: 42},
			wantErr: nil,
		},
		{
			name:    "empty JSON object",
			body:    `{}`,
			want:    &TestData{},
			wantErr: nil,
		},
		{
			name:    "invalid JSON body",
			body:    `{invalid json`,
			want:    nil,
			wantErr: ErrInternal,
		},
		{
			name:    "empty body",
			body:    "",
			want:    nil,
			wantErr: ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			got, err := ReadBody[TestData](req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ReadBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && (got.Name != tt.want.Name || got.Value != tt.want.Value) {
				t.Errorf("ReadBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseMessage(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{
			name:    "200 OK with message",
			code:    http.StatusOK,
			message: "success",
		},
		{
			name:    "201 Created with message",
			code:    http.StatusCreated,
			message: "resource created",
		},
		{
			name:    "empty message",
			code:    http.StatusOK,
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ResponseMessage(w, tt.code, tt.message)

			if w.Code != tt.code {
				t.Errorf("ResponseMessage() status = %v, want %v", w.Code, tt.code)
			}

			var resp StandardResponse[any]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Message == nil && tt.message != "" {
				t.Error("ResponseMessage() message is nil, want non-nil")
			} else if resp.Message != nil && *resp.Message != tt.message {
				t.Errorf("ResponseMessage() message = %v, want %v", *resp.Message, tt.message)
			}

			if resp.Error != nil {
				t.Error("ResponseMessage() error should be nil")
			}
			if resp.Data != nil {
				t.Error("ResponseMessage() data should be nil")
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		errMsg  string
	}{
		{
			name:   "400 Bad Request",
			code:   http.StatusBadRequest,
			errMsg: "invalid input",
		},
		{
			name:   "500 Internal Server Error",
			code:   http.StatusInternalServerError,
			errMsg: "internal error",
		},
		{
			name:   "empty error message",
			code:   http.StatusBadRequest,
			errMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ResponseError(w, tt.code, tt.errMsg)

			if w.Code != tt.code {
				t.Errorf("ResponseError() status = %v, want %v", w.Code, tt.code)
			}

			var resp StandardResponse[any]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Error == nil && tt.errMsg != "" {
				t.Error("ResponseError() error is nil, want non-nil")
			} else if resp.Error != nil && *resp.Error != tt.errMsg {
				t.Errorf("ResponseError() error = %v, want %v", *resp.Error, tt.errMsg)
			}

			if resp.Message != nil {
				t.Error("ResponseError() message should be nil")
			}
			if resp.Data != nil {
				t.Error("ResponseError() data should be nil")
			}
		})
	}
}

func TestResponseJSON(t *testing.T) {
	type TestData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name string
		code int
		data TestData
	}{
		{
			name: "200 OK with data",
			code: http.StatusOK,
			data: TestData{ID: 1, Name: "test"},
		},
		{
			name: "201 Created with data",
			code: http.StatusCreated,
			data: TestData{ID: 2, Name: "created"},
		},
		{
			name: "empty data",
			code: http.StatusOK,
			data: TestData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ResponseJSON(w, tt.code, tt.data)

			if w.Code != tt.code {
				t.Errorf("ResponseJSON() status = %v, want %v", w.Code, tt.code)
			}

			var resp StandardResponse[TestData]
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Data == nil {
				t.Error("ResponseJSON() data is nil, want non-nil")
			} else if resp.Data.ID != tt.data.ID || resp.Data.Name != tt.data.Name {
				t.Errorf("ResponseJSON() data = %v, want %v", resp.Data, tt.data)
			}

			if resp.Message != nil {
				t.Error("ResponseJSON() message should be nil")
			}
			if resp.Error != nil {
				t.Error("ResponseJSON() error should be nil")
			}
		})
	}
}

func TestStandardResponseStructure(t *testing.T) {
	t.Run("response with all fields nil", func(t *testing.T) {
		resp := StandardResponse[any]{}
		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(data, &result)

		// All fields should be null in JSON
		if result["message"] != nil {
			t.Error("message should be null")
		}
		if result["data"] != nil {
			t.Error("data should be null")
		}
		if result["error"] != nil {
			t.Error("error should be null")
		}
	})
}
