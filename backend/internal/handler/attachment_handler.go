package handler

import (
	"net/http"

	"livekit-consulting/backend/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AttachmentHandler struct {
	fileStorage service.FileStorage
}

func NewAttachmentHandler(fileStorage service.FileStorage) *AttachmentHandler {
	return &AttachmentHandler{fileStorage: fileStorage}
}

func (h *AttachmentHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := uuid.Parse(vars["roomId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error retrieving the file")
		return
	}
	defer file.Close()

	attachment, err := h.fileStorage.UploadFile(r.Context(), handler, roomID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upload file")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"fileId": attachment.ID.String()})
}
