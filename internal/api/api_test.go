package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBreakdownHandler(t *testing.T) {
	s := NewServer("example")
	body := `{"p":760,"d":1,"A":15,"B":365,"gamma":0.01}`
	req := httptest.NewRequest(http.MethodPost, "/api/breakdown", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/breakdown status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("content-type = %q, want application/json", ct)
	}
	if !strings.Contains(rec.Body.String(), "voltage") {
		t.Errorf("response missing voltage field: %s", rec.Body.String())
	}
}

func TestBreakdownRejectsBadInput(t *testing.T) {
	s := NewServer("example")
	req := httptest.NewRequest(http.MethodPost, "/api/breakdown", strings.NewReader(`{"p":-1,"d":1}`))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestSweepHandler(t *testing.T) {
	s := NewServer("example")
	body := `{"p":760,"d_min":0.05,"d_max":20,"points":50,"log":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/sweep", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/sweep status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "points") {
		t.Errorf("response missing points field: %s", rec.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := NewServer("example")
	req := httptest.NewRequest(http.MethodGet, "/api/breakdown", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET should be 405, got %d", rec.Code)
	}
}
