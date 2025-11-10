package handler

import (
	"encoding/json"
	"net/http"

	"livekit-consulting/backend/internal/service"
)

type AgentWebhookHandler struct {
	messageService *service.MessageService
}

func NewAgentWebhookHandler(messageService *service.MessageService) *AgentWebhookHandler {
	return &AgentWebhookHandler{
		messageService: messageService,
	}
}

func (h *AgentWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload service.AgentWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate payload if necessary (e.g., check event type, room name)
	if payload.Event != "transcript_uploaded" {
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	_, err := h.messageService.CreateTranscriptMessage(r.Context(), &payload)
	if err != nil {
		http.Error(w, "Failed to create transcript message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Transcript webhook processed"})
}
