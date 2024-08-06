package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/claraferraz/ask-me-anything/internal/store/pgstore"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

func SendJSON(w http.ResponseWriter, response any) {
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (h apiHandler) ReadRoom(w http.ResponseWriter, r *http.Request) (room pgstore.Room, rawRoomID string, roomID uuid.UUID, ok bool) {
	rawRoomID = chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	room, err = h.q.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	return room, rawRoomID, roomID, true

}
