// Package server wires HTTP handlers and runs the HTTP server.
package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"git.skobk.in/skobkin/ip-detect/internal/clientinfo"
	"git.skobk.in/skobkin/ip-detect/internal/config"
	"git.skobk.in/skobkin/ip-detect/internal/templates"
)

type handler struct {
	cfg    config.Config
	logger *slog.Logger
	tpl    *template.Template
}

type viewModel struct {
	Data      clientinfo.Data
	JSONPath  string
	PlainPath string
	Timestamp string
}

func newHandler(cfg config.Config, logger *slog.Logger) (http.Handler, error) {
	tpl, err := templates.Client()
	if err != nil {
		return nil, fmt.Errorf("load template: %w", err)
	}

	return &handler{cfg: cfg, logger: logger, tpl: tpl}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	data := clientinfo.Collect(r.Context(), r, h.cfg)
	lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

	switch r.URL.Path {
	case "/", "":
		h.respondHTML(lrw, data)
	case "/json":
		h.respondJSON(lrw, data)
	case "/plain":
		h.respondPlain(lrw, data)
	default:
		http.NotFound(lrw, r)
	}

	h.logger.Info("request completed",
		"method", r.Method,
		"path", r.URL.Path,
		"status", lrw.status,
		"ip", data.IPAddress,
		"hostname", hostnameValue(data.Hostname),
		"duration", time.Since(start),
	)
}

func (h *handler) respondHTML(w http.ResponseWriter, data clientinfo.Data) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	model := viewModel{
		Data:      data,
		JSONPath:  "/json",
		PlainPath: "/plain",
	}

	if data.Timestamp != nil {
		model.Timestamp = data.Timestamp.Format(time.DateTime)
	}

	if err := h.tpl.Execute(w, model); err != nil {
		h.logger.Error("template render failed", "error", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *handler) respondJSON(w http.ResponseWriter, data clientinfo.Data) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("json response failed", "error", err)
	}
}

func (h *handler) respondPlain(w http.ResponseWriter, data clientinfo.Data) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := fmt.Fprintf(w, "%s\n", data.IPAddress); err != nil {
		h.logger.Error("plaintext response failed", "error", err)
	}
}

func hostnameValue(hostname *string) string {
	if hostname == nil {
		return ""
	}

	return *hostname
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err != nil {
		return n, fmt.Errorf("write response: %w", err)
	}

	return n, nil
}
