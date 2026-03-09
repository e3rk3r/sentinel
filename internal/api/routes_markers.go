package api

import "net/http"

func (h *Handler) registerMarkersRoutes(mux *http.ServeMux) {
	h.registerRoutes(mux, []routeBinding{
		{pattern: "GET /api/ops/markers", handler: h.listMarkerPatterns},
		{pattern: "PUT /api/ops/markers/{pattern}", handler: h.upsertMarkerPattern},
		{pattern: "DELETE /api/ops/markers/{pattern}", handler: h.deleteMarkerPattern},
	})
}
