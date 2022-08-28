package handler

import (
	"github.com/rs/zerolog"
	"net/http"
)

type PACFileHandler struct {
	logger   zerolog.Logger
	filePath string
}

func NewPACFileHandler(filePath string, logger zerolog.Logger) *PACFileHandler {
	return &PACFileHandler{
		logger:   logger,
		filePath: filePath,
	}
}

func (h *PACFileHandler) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	http.ServeFile(w, r, h.filePath)
}
