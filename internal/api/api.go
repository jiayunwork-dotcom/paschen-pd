// Package api exposes the Paschen breakdown model over HTTP. It serves the
// static web console from web/ and a JSON endpoint:
//
//	POST /api/breakdown   compute the breakdown voltage for one pressure+distance
//	POST /api/sweep       sample the breakdown voltage over a pd range
//
// All computation is delegated to the paschen and sweep packages; this package
// only decodes requests, calls the engine, and encodes results or errors.
// Invalid input yields a JSON error body of the form {"error": "..."} with a
// 4xx status so the failure is visible to both API clients and the web page.
package api

import (
	"encoding/json"
	"net/http"

	"paschen-pd/internal/paschen"
	"paschen-pd/internal/sweep"
)

// Server bundles the HTTP routing for paschen-pd. It serves the static web
// console from WebDir and the example files from ExampleDir.
type Server struct {
	ExampleDir string
	WebDir     string
	mux        *http.ServeMux
}

// NewServer constructs a Server with the web console in "web" and example
// files in dir.
func NewServer(dir string) *Server {
	s := &Server{
		ExampleDir: dir,
		WebDir:     "web",
		mux:        http.NewServeMux(),
	}
	s.routes()
	return s
}

// Handler returns the underlying http.Handler (used by tests via httptest).
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the HTTP server on addr with a short read-header
// timeout so a stalled client cannot tie up the server.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

// routes wires the endpoints and static file serving. API paths are registered
// first so they take precedence over the catch-all root handler.
//
// The brief names the endpoints /api/vs (single breakdown) and /api/curve
// (pd sweep); they are exposed alongside the internal /api/breakdown and
// /api/sweep names so both the web console and the spec map to the same engine.
func (s *Server) routes() {
	s.mux.HandleFunc("/api/breakdown", s.handleBreakdown)
	s.mux.HandleFunc("/api/vs", s.handleBreakdown)
	s.mux.HandleFunc("/api/sweep", s.handleSweep)
	s.mux.HandleFunc("/api/curve", s.handleSweep)
	s.mux.Handle(
		"/example/",
		http.StripPrefix("/example/", http.FileServer(http.Dir(s.ExampleDir))),
	)
	s.mux.Handle("/", http.FileServer(http.Dir(s.WebDir)))
}

func (s *Server) handleBreakdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req BreakdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := ValidateBreakdown(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	params := paramsFrom(req)
	pd := req.P * req.D / 10.0 // Torr·mm -> Torr·cm
	v := params.BreakdownVoltage(pd)
	pdMin, vMin := params.Minimum()
	resp := BreakdownResponse{
		PD:            pd,
		Voltage:       v,
		PDMin:         pdMin,
		VMin:          vMin,
		A:             params.A,
		B:             params.B,
		Gamma:         params.Gamma,
		Breakdownable: v > 0,
		LeftBranch:    pd < pdMin,
	}
	writeJSON(w, resp)
}

// SweepRequest is documented in response.go; the handler follows.

func (s *Server) handleSweep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req SweepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := ValidateSweep(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a, b, g := defaultIfZero(BreakdownRequest{A: req.A, B: req.B, Gamma: req.Gamma})
	params := paschen.Params{A: a, B: b, Gamma: g}
	cfg := sweep.Config{
		From:   req.P * req.DMin / 10.0,
		To:     req.P * req.DMax / 10.0,
		Points: req.Points,
		Log:    req.Log,
		Params: params,
	}
	raw := sweep.Run(cfg)
	pts := make([]SweepPoint, 0, len(raw))
	for _, p := range raw {
		pts = append(pts, SweepPoint{PD: p.PD, Voltage: p.Voltage})
	}
	writeJSON(w, SweepResponse{Points: pts, A: a, B: b, Gamma: g})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
