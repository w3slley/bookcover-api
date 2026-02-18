package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccess_StatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	Success(rr, "https://example.com/cover.jpg")

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestSuccess_JSONBody(t *testing.T) {
	rr := httptest.NewRecorder()
	body := Success(rr, "https://example.com/cover.jpg")

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if result["url"] != "https://example.com/cover.jpg" {
		t.Errorf("expected url field to be the given URL, got %q", result["url"])
	}
}

func TestSuccess_NoHTMLEscaping(t *testing.T) {
	rr := httptest.NewRecorder()
	// URLs with & should not be escaped to \u0026
	body := Success(rr, "https://example.com/img?size=large&quality=high")

	if strings.Contains(string(body), `\u0026`) {
		t.Errorf("expected & not to be HTML-escaped, but got: %s", string(body))
	}
	if !strings.Contains(string(body), "&") {
		t.Errorf("expected & to be preserved in URL, got: %s", string(body))
	}
}

func TestSuccess_EmptyURL(t *testing.T) {
	rr := httptest.NewRecorder()
	body := Success(rr, "")

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if result["url"] != "" {
		t.Errorf("expected empty url field, got %q", result["url"])
	}
}

func TestError_StatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"not found", http.StatusNotFound},
		{"bad request", http.StatusBadRequest},
		{"internal server error", http.StatusInternalServerError},
		{"unauthorized", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			Error(rr, tt.statusCode, "some error")

			if rr.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}

func TestError_JSONBody(t *testing.T) {
	rr := httptest.NewRecorder()
	body := Error(rr, http.StatusNotFound, "book not found")

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if result["error"] != "book not found" {
		t.Errorf("expected error field to be 'book not found', got %q", result["error"])
	}
}

func TestError_EmptyMessage(t *testing.T) {
	rr := httptest.NewRecorder()
	body := Error(rr, http.StatusBadRequest, "")

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if result["error"] != "" {
		t.Errorf("expected empty error field, got %q", result["error"])
	}
}
