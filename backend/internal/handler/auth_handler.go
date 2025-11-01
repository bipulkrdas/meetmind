package handler

import (
    "encoding/json"
    "net/http"
    "livekit-consulting/backend/internal/model"
    "livekit-consulting/backend/internal/service"
    "livekit-consulting/backend/internal/utils"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    var req model.UserSignUpRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    err := h.authService.SignUp(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, map[string]string{
        "message": "User created successfully",
    })
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
    var req model.UserSignInRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    authResp, err := h.authService.SignIn(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
    var req model.PasswordResetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    h.authService.RequestPasswordReset(r.Context(), req.Email)

    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "If the email exists, a reset link has been sent",
    })
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req model.PasswordResetConfirm
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    err := h.authService.ResetPassword(r.Context(), &req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Password reset successfully",
    })
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
    userID, err := getUserIDFromContext(r)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    user, err := h.authService.GetMe(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, user)
}
