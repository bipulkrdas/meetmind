package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/service"
	"livekit-consulting/backend/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid room ID")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req model.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	message, err := h.messageService.CreateMessage(r.Context(), &req, roomID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid room ID")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	var before *uuid.UUID
	if b := r.URL.Query().Get("before"); b != "" {
		beforeID, err := uuid.Parse(b)
		if err == nil {
			before = &beforeID
		}
	}

	messages, err := h.messageService.GetMessages(r.Context(), roomID, userID, limit, before)
	if err != nil {
		log.Error().
			Err(err).
			Str("room_id", roomID.String()).
			Str("user_id", userID.String()).
			Msg("Failed to get messages")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, err := uuid.Parse(vars["messageId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req model.UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.messageService.UpdateMessage(r.Context(), messageID, userID, req.Content); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "message updated successfully"})
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, err := uuid.Parse(vars["messageId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.messageService.DeleteMessage(r.Context(), messageID, userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "message deleted successfully"})
}

func (h *MessageHandler) UpdateLastRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid room ID")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req model.UpdateLastReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.messageService.UpdateLastRead(r.Context(), roomID, userID, req.LastReadSequenceNumber)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "last read updated"})
}
