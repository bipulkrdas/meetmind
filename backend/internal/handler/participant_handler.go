package handler

import (
    "encoding/json"
    "net/http"
    "livekit-consulting/backend/internal/model"
    "livekit-consulting/backend/internal/service"
    "livekit-consulting/backend/internal/utils"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

type ParticipantHandler struct {
    participantService *service.ParticipantService
}

func NewParticipantHandler(participantService *service.ParticipantService) *ParticipantHandler {
    return &ParticipantHandler{participantService: participantService}
}

func (h *ParticipantHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
    inviterID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    var req model.AddParticipantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    participant, err := h.participantService.AddParticipant(r.Context(), roomID, inviterID, &req)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, participant)
}

func (h *ParticipantHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    participants, err := h.participantService.GetRoomParticipants(r.Context(), roomID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, participants)
}

func (h *ParticipantHandler) RemoveParticipant(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    participantID, err := uuid.Parse(vars["participantId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid participant ID")
        return
    }

    err = h.participantService.RemoveParticipant(r.Context(), roomID, participantID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Participant removed successfully"})
}

func (h *ParticipantHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    token := r.URL.Query().Get("token")
    if token == "" {
        respondWithError(w, http.StatusBadRequest, "Missing join token")
        return
    }

    livekitToken, err := h.participantService.GenerateParticipantToken(r.Context(), roomID, token)
    if err != nil {
        respondWithError(w, http.StatusForbidden, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"livekit_token": livekitToken})
}

func (h *ParticipantHandler) JoinRoomInternal(w http.ResponseWriter, r *http.Request) {
    userID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    livekitToken, err := h.participantService.GenerateInternalParticipantToken(r.Context(), roomID, userID)
    if err != nil {
        if err.Error() == "access denied: user is not a participant of this room" {
            respondWithError(w, http.StatusForbidden, err.Error())
        } else {
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"livekit_token": livekitToken})
}

func (h *ParticipantHandler) InviteParticipantsToJoinMeeting(w http.ResponseWriter, r *http.Request) {
	inviterID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	err = h.participantService.InviteParticipantsToJoinMeeting(r.Context(), roomID, inviterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Invitations sent successfully"})
}

func (h *ParticipantHandler) GenerateMeetingUrl(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	meetingURL, err := h.participantService.GenerateMeetingUrl(r.Context(), roomID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"meeting_url": meetingURL})
}
