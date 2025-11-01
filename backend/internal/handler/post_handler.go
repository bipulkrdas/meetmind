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

type PostHandler struct {
    postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
    return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

    var req model.CreatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    post, err := h.postService.CreatePost(r.Context(), userID, roomID, &req)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomID, err := uuid.Parse(vars["roomId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid room ID")
        return
    }

    posts, err := h.postService.GetPosts(r.Context(), roomID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
    userID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    vars := mux.Vars(r)
    postID, err := uuid.Parse(vars["postId"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    err = h.postService.DeletePost(r.Context(), postID, userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}
