package service

import (
    "context"
    "errors"
    "time"
    "livekit-consulting/backend/internal/model"
    "livekit-consulting/backend/internal/repository"
    "livekit-consulting/backend/internal/utils"
    "livekit-consulting/backend/internal/service/email"

    "github.com/google/uuid"
)

type AuthService struct {
    userRepo       repository.UserRepository
    resetTokenRepo repository.PasswordResetTokenRepository
    emailService   *email.EmailService
    livekitService *LiveKitService
    jwtSecret      string
    frontendURL    string
}

func NewAuthService(
    userRepo repository.UserRepository,
    resetTokenRepo repository.PasswordResetTokenRepository,
    emailService *email.EmailService,
    livekitService *LiveKitService,
    jwtSecret, frontendURL string,
) *AuthService {
    return &AuthService{
        userRepo:       userRepo,
        resetTokenRepo: resetTokenRepo,
        emailService:   emailService,
        livekitService: livekitService,
        jwtSecret:      jwtSecret,
        frontendURL:    frontendURL,
    }
}

func (s *AuthService) SignUp(ctx context.Context, req *model.UserSignUpRequest) error {
    if req.Password != req.ConfirmPassword {
        return errors.New("passwords do not match")
    }
    
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return errors.New("email already exists")
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return err
    }
    
    user := &model.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Name:         req.Name,
    }
    
    return s.userRepo.Create(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, req *model.UserSignInRequest) (*model.AuthResponse, error) {
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }
    
    token, expiresAt, err := utils.GenerateJWT(user.ID.String(), user.Email, s.jwtSecret)
    if err != nil {
        return nil, err
    }
    
    livekitToken, err := s.livekitService.GenerateToken(
        user.ID.String(),
        "default",
        true,      
        true,      
    )
    if err != nil {
        return nil, err
    }
    
    s.userRepo.UpdateLastLogin(ctx, user.ID)
    
    return &model.AuthResponse{
        Token:        token,
        User:         *user,
        LiveKitToken: livekitToken,
        ExpiresAt:    expiresAt,
    }, nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil
    }
    
    resetToken := uuid.New().String()
    expiresAt := time.Now().Add(1 * time.Hour)
    
    err = s.resetTokenRepo.Create(ctx, user.ID, resetToken, expiresAt)
    if err != nil {
        return err
    }
    
    resetURL := s.frontendURL + "/auth/reset-password-form?token=" + resetToken
    return s.emailService.SendPasswordResetEmail(ctx, user.Email, resetToken, resetURL)
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*model.User, error) {
    return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) ResetPassword(ctx context.Context, req *model.PasswordResetConfirm) error {
    resetToken, err := s.resetTokenRepo.GetByToken(ctx, req.Token)
    if err != nil {
        return errors.New("invalid or expired token")
    }
    
    if resetToken.Used || time.Now().After(resetToken.ExpiresAt) {
        return errors.New("invalid or expired token")
    }
    
    user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
    if err != nil {
        return err
    }
    
    hashedPassword, err := utils.HashPassword(req.NewPassword)
    if err != nil {
        return err
    }
    
    user.PasswordHash = hashedPassword
    err = s.userRepo.Update(ctx, user)
    if err != nil {
        return err
    }
    
    return s.resetTokenRepo.MarkAsUsed(ctx, resetToken.ID)
}
