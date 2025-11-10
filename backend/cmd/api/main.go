package main

import (
	"context"
	"net/http"
	"os"

	"livekit-consulting/backend/internal/config"
	"livekit-consulting/backend/internal/database"
	"livekit-consulting/backend/internal/handler"
	"livekit-consulting/backend/internal/middleware"
	"livekit-consulting/backend/internal/repository"
	"livekit-consulting/backend/internal/service"
	"livekit-consulting/backend/internal/service/email"
	"livekit-consulting/backend/pkg/logger"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	postRepo := repository.NewPostRepository(db)
	resetTokenRepo := repository.NewPasswordResetTokenRepository(db)
	inviteRepo := repository.NewInviteRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	var emailProvider email.EmailProvider
	if cfg.EmailProvider == "sendgrid" {
		emailProvider = email.NewSendGridProvider(
			cfg.SendGridAPIKey,
			cfg.SendGridFromEmail,
			cfg.SendGridFromName,
		)
	} else {
		emailProvider = email.NewMailjetProvider(
			cfg.MailjetAPIKey,
			cfg.MailjetSecretKey,
			cfg.MailjetFromEmail,
			cfg.MailjetFromName,
		)
	}

	emailService := email.NewEmailService(
		emailProvider,
		cfg.FromEmail,
		cfg.FromName,
	)

	livekitService := service.NewLiveKitService(
		cfg.LiveKitAPIKey,
		cfg.LiveKitAPISecret,
		cfg.LiveKitURL,
	)

	authService := service.NewAuthService(
		userRepo,
		resetTokenRepo,
		emailService,
		livekitService,
		cfg.JWTSecret,
		cfg.FrontendURL,
	)

	roomService := service.NewRoomService(
		roomRepo,
		participantRepo,
		livekitService,
	)

	participantService := service.NewParticipantService(
		participantRepo,
		roomRepo,
		inviteRepo,
		emailService,
		livekitService,
		cfg.FrontendURL,
	)

	postService := service.NewPostService(postRepo, roomRepo)

	var fileStorage service.FileStorage
	switch cfg.StorageProvider {
	case "minio":
		fileStorage, err = service.NewMinioFileStorage(cfg, attachmentRepo)
	case "s3":
		fileStorage, err = service.NewS3FileStorage(cfg, attachmentRepo)
	case "gcs":
		fileStorage, err = service.NewGCSFileStorage(context.Background(), cfg, attachmentRepo)
	default:
		log.Fatal().Msg("Unsupported storage provider")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create file storage")
	}

	messageService := service.NewMessageService(messageRepo, participantRepo, attachmentRepo, roomRepo)

	s3TranscriptStorage, err := service.NewS3TranscriptStorage(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create S3 transcript storage")
	}

	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService)
	participantHandler := handler.NewParticipantHandler(participantService)
	postHandler := handler.NewPostHandler(postService)
	messageHandler := handler.NewMessageHandler(messageService)
	attachmentHandler := handler.NewAttachmentHandler(fileStorage)
	agentWebhookHandler := handler.NewAgentWebhookHandler(messageService)
	transcriptHandler := handler.NewTranscriptHandler(messageRepo, s3TranscriptStorage)

	r := mux.NewRouter()

	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))
	r.Use(middleware.LoggingMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/agent-webhook", agentWebhookHandler.HandleWebhook).Methods("POST")

	api.HandleFunc("/auth/signup", authHandler.SignUp).Methods("POST")
	api.HandleFunc("/auth/signin", authHandler.SignIn).Methods("POST")
	api.HandleFunc("/auth/reset-password", authHandler.RequestPasswordReset).Methods("POST")
	api.HandleFunc("/auth/reset-password/confirm", authHandler.ResetPassword).Methods("POST")
	api.HandleFunc("/rooms/{roomId}/join_external", participantHandler.JoinRoom).Methods("POST")

	authAPI := api.PathPrefix("/app").Subrouter()
	authAPI.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))

	authAPI.HandleFunc("/rooms/{roomId}/transcript/{messageId}/{s3KeyPath:.+}", transcriptHandler.GetTranscript).Methods("GET")
	authAPI.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	authAPI.HandleFunc("/rooms", roomHandler.GetUserRooms).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}", roomHandler.GetRoomDetails).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}/livekit_create", roomHandler.CreateRoomAtLiveKit).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}", roomHandler.DeleteRoom).Methods("DELETE")

	authAPI.HandleFunc("/rooms/{roomId}/participants", participantHandler.AddParticipant).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/participants", participantHandler.GetParticipants).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}/participants/{participantId}", participantHandler.RemoveParticipant).Methods("DELETE")
	// authAPI.HandleFunc("/rooms/{roomId}/join_external", participantHandler.JoinRoom).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/join_internal", participantHandler.JoinRoomInternal).Methods("POST")

	authAPI.HandleFunc("/rooms/{roomId}/invite_participants_to_join_meeting", participantHandler.InviteParticipantsToJoinMeeting).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/generate_meeting_url", participantHandler.GenerateMeetingUrl).Methods("POST")

	authAPI.HandleFunc("/rooms/{roomId}/posts", postHandler.CreatePost).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/posts", postHandler.GetPosts).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}/posts/{postId}", postHandler.DeletePost).Methods("DELETE")

	authAPI.HandleFunc("/rooms/{roomId}/messages", messageHandler.CreateMessage).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/messages", messageHandler.GetMessages).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}/update_last_read_for_user", messageHandler.UpdateLastRead).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/attachments", attachmentHandler.UploadAttachment).Methods("POST")

	authAPI.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
