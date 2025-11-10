package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/repository"
	"livekit-consulting/backend/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TranscriptHandler struct {
	messageRepo         repository.MessageRepository
	s3TranscriptStorage service.S3TranscriptStorage
}

func NewTranscriptHandler(messageRepo repository.MessageRepository, s3TranscriptStorage service.S3TranscriptStorage) *TranscriptHandler {
	return &TranscriptHandler{
		messageRepo:         messageRepo,
		s3TranscriptStorage: s3TranscriptStorage,
	}
}

func (h *TranscriptHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomIDStr := vars["roomId"]
	messageIDStr := vars["messageId"]
	s3KeyPath := vars["s3KeyPath"] // This will be "s3_keys.json" or similar

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Fetch the message to verify it's a transcript message and get the actual S3 key
	message, err := h.messageRepo.GetByID(r.Context(), messageID)
	if err != nil {
		http.Error(w, "Message not found or database error", http.StatusInternalServerError)
		return
	}
	if message == nil || message.RoomID != roomID || message.MessageType != model.MessageTypeMeetingTranscript || message.ExtraData == nil || message.ExtraData.Transcript == nil {
		http.Error(w, "Transcript message not found or invalid", http.StatusNotFound)
		return
	}

	// Determine the correct S3 key based on s3KeyPath
	var actualS3Key string
	if s3KeyPath == "s3_keys.json" { // Assuming the frontend sends this literal string
		actualS3Key = message.ExtraData.Transcript.S3Keys.JSON
	} else if s3KeyPath == "s3_keys.text" { // If we ever want to support text directly
		actualS3Key = message.ExtraData.Transcript.S3Keys.Text
	} else {
		http.Error(w, "Invalid S3 key path requested", http.StatusBadRequest)
		return
	}

	if actualS3Key == "" {
		http.Error(w, "S3 key not found in message extra data", http.StatusNotFound)
		return
	}

	// Get the transcript file from S3
	fileReader, err := h.s3TranscriptStorage.GetTranscriptFile(r.Context(), actualS3Key)
	if err != nil {
		http.Error(w, "Failed to retrieve transcript from storage: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileReader.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, fileReader); err != nil {
		// Log the error, but response might already be partially sent
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to stream transcript content"})
	}
}
