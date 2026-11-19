package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApiAsciiArtHandlerSuccess(t *testing.T) {
	rec := postJSON(t, `{"text":"Hi","banner":"standard"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var res ApiResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Result == "" {
		t.Fatal("expected generated ASCII result")
	}
}

func TestApiAsciiArtHandlerRejectsUnsafeRequests(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
	}{
		{
			name:        "wrong content type",
			body:        `{"text":"Hi","banner":"standard"}`,
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name:        "unknown field",
			body:        `{"text":"Hi","banner":"standard","extra":true}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "unsupported character",
			body:        `{"text":"¿","banner":"standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "invalid banner",
			body:        `{"text":"Hi","banner":"../../standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "text too long",
			body:        `{"text":"` + strings.Repeat("a", maxTextBytes+1) + `","banner":"standard"}`,
			contentType: "application/json",
			wantStatus:  http.StatusRequestEntityTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()

			ApiAsciiArtHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestApiAsciiArtHandlerRejectsOversizedBody(t *testing.T) {
	body := `{"text":"` + strings.Repeat("a", maxRequestBytes) + `","banner":"standard"}`
	rec := postJSON(t, body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
}

func TestApiAsciiArtHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ascii-art", nil)
	rec := httptest.NewRecorder()

	ApiAsciiArtHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if rec.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("Allow = %q, want %q", rec.Header().Get("Allow"), http.MethodPost)
	}
}

func TestHomeAndStaticHandlers(t *testing.T) {
	t.Run("home renders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		HomeHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Body.String(), "Banner Generator") {
			t.Fatal("expected home page content")
		}
	})

	t.Run("unknown page is not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rec := httptest.NewRecorder()

		HomeHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("static directory is not listed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/static/", nil)
		rec := httptest.NewRecorder()

		StaticHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rec, req)

	requiredHeaders := []string{
		"Content-Security-Policy",
		"Permissions-Policy",
		"Referrer-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
	}
	for _, header := range requiredHeaders {
		if rec.Header().Get(header) == "" {
			t.Fatalf("expected %s header", header)
		}
	}
}

func postJSON(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ApiAsciiArtHandler(rec, req)
	return rec
}

func TestExportHandler(t *testing.T) {
	t.Run("GET txt export", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/export?text=Hello&banner=standard&format=txt", nil)
		rec := httptest.NewRecorder()

		ExportHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if contentType := rec.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
			t.Fatalf("Content-Type = %q, want %q", contentType, "text/plain; charset=utf-8")
		}
		if contentDisp := rec.Header().Get("Content-Disposition"); !strings.Contains(contentDisp, "attachment; filename=\"ascii-art-standard.txt\"") {
			t.Fatalf("Content-Disposition = %q, want filename containing standard.txt", contentDisp)
		}
		contentLength := rec.Header().Get("Content-Length")
		if contentLength == "" {
			t.Fatal("expected Content-Length header")
		}
		if len(rec.Body.Bytes()) == 0 {
			t.Fatal("expected non-empty body")
		}
	})

	t.Run("POST html export", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader("text=Hello&banner=shadow&format=html"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		ExportHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		if contentType := rec.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
			t.Fatalf("Content-Type = %q, want %q", contentType, "text/html; charset=utf-8")
		}
		if contentDisp := rec.Header().Get("Content-Disposition"); !strings.Contains(contentDisp, "attachment; filename=\"ascii-art-shadow.html\"") {
			t.Fatalf("Content-Disposition = %q, want filename containing shadow.html", contentDisp)
		}
		if !strings.Contains(rec.Body.String(), "<html lang=\"en\">") {
			t.Fatal("expected HTML format content")
		}
	})

	t.Run("POST json export", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/export", strings.NewReader("text=Hello&banner=standard&format=json"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		ExportHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if contentType := rec.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
			t.Fatalf("Content-Type = %q, want %q", contentType, "application/json; charset=utf-8")
		}
		if contentDisp := rec.Header().Get("Content-Disposition"); !strings.Contains(contentDisp, "attachment; filename=\"ascii-art-standard.json\"") {
			t.Fatalf("Content-Disposition = %q, want filename containing standard.json", contentDisp)
		}
		var parsed map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
			t.Fatalf("failed to decode JSON response: %v", err)
		}
		if parsed["ascii"] == "" {
			t.Fatal("expected ascii key in JSON output")
		}
	})

	t.Run("invalid inputs", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			url        string
			body       string
			wantStatus int
		}{
			{
				name:       "empty text",
				method:     http.MethodGet,
				url:        "/export?text=&banner=standard&format=txt",
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "invalid banner",
				method:     http.MethodGet,
				url:        "/export?text=Hello&banner=invalid&format=txt",
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "unsupported character",
				method:     http.MethodGet,
				url:        "/export?text=Hello¿&banner=standard&format=txt",
				wantStatus: http.StatusBadRequest,
			},
			{
				name:       "text too long",
				method:     http.MethodGet,
				url:        "/export?text=" + strings.Repeat("a", maxTextBytes+1) + "&banner=standard&format=txt",
				wantStatus: http.StatusRequestEntityTooLarge,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(tt.method, tt.url, nil)
				rec := httptest.NewRecorder()

				ExportHandler(rec, req)

				if rec.Code != tt.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
				}
			})
		}
	})
}
