package http

import (
	"encoding/json"
	"imageResample/pkg/api"
	"imageResample/pkg/utils"
	"log/slog"
	"net/http"
)

type ImageResampler interface {
	Resample(request api.ImageRequest) (int64, bool, error)
}

type ResizeHandler struct {
	service ImageResampler
	log     *slog.Logger
}

func NewResizeHandler(service ImageResampler, log *slog.Logger) *ResizeHandler {
	return &ResizeHandler{service: service, log: log}
}

func (h *ResizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Received request", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	var req api.ImageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to parse request", slog.Any("error", err))
		utils.JSON(w, http.StatusBadRequest, api.ImageErrorResponse{Error: err.Error()})
		return
	}

	processingTime, cached, err := h.service.Resample(req)
	if err != nil {
		h.log.Error(err.Error())
		utils.JSON(w, http.StatusBadRequest, api.ImageErrorResponse{Error: err.Error()})
		return
	}

	response := api.ImageSuccessResponse{
		ProcessingTime: processingTime,
		Cached:         cached,
	}

	utils.JSON(w, http.StatusOK, response)
	return
}
