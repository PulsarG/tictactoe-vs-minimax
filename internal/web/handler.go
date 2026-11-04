package web

import (
	"encoding/json"
	"net/http"
	"tictactoe-vs-minmax/internal"
)

type GameHandler struct {
	svc internal.GameService
}

func NewGameHandler(svc internal.GameService) *GameHandler {
	return &GameHandler{svc: svc}
}

func (h *GameHandler) PostGame(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("uuid")
	id, err := ParseUUID(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	game := internal.Game{ID: id, Board: ToDomainBoard(req)}
	update, err := h.svc.NextMove(game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ToResponse(update)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GameHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /game/{uuid}", h.PostGame)
}
