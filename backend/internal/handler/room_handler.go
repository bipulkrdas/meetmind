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

type RoomHandler struct {
    roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
    return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
    userID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    var req model.CreateRoomRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    room, err := h.roomService.CreateRoom(r.Context(), userID, &req)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, room)
}

func (h *RoomHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
    userID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    rooms, err := h.roomService.GetUserRooms(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoomDetails(w http.ResponseWriter, r *http.Request) {
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

    room, err := h.roomService.GetRoomDetails(r.Context(), roomID, userID)
    if err != nil {
        respondWithError(w, http.StatusForbidden, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
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

    err = h.roomService.DeleteRoom(r.Context(), roomID, userID)
    if err != nil {
        respondWithError(w, http.StatusForbidden, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Room deleted successfully",
    })
}

func (h *RoomHandler) CreateRoomAtLiveKit(w http.ResponseWriter, r *http.Request) {
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

    lkRoom, err := h.roomService.CreateLiveKitRoom(r.Context(), roomID, userID)
    if err != nil {
        if err.Error() == "permission denied: only room owner can create livekit room" {
            respondWithError(w, http.StatusForbidden, err.Error())
        } else {
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusCreated, lkRoom)
}
