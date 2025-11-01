package main

import (
	"log"
	"net/http"
	"os"

	"livekit-consulting/backend/internal/config"
	"livekit-consulting/backend/internal/database"
	"livekit-consulting/backend/internal/handler"
	"livekit-consulting/backend/internal/middleware"
	"livekit-consulting/backend/internal/repository"
	"livekit-consulting/backend/internal/service"
	"livekit-consulting/backend/internal/service/email"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	postRepo := repository.NewPostRepository(db)
	resetTokenRepo := repository.NewPasswordResetTokenRepository(db)
	inviteRepo := repository.NewInviteRepository(db)

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

	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService)
	participantHandler := handler.NewParticipantHandler(participantService)
	postHandler := handler.NewPostHandler(postService)

	r := mux.NewRouter()

	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))
	r.Use(middleware.LoggingMiddleware)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/auth/signup", authHandler.SignUp).Methods("POST")
	api.HandleFunc("/auth/signin", authHandler.SignIn).Methods("POST")
	api.HandleFunc("/auth/reset-password", authHandler.RequestPasswordReset).Methods("POST")
	api.HandleFunc("/auth/reset-password/confirm", authHandler.ResetPassword).Methods("POST")
	api.HandleFunc("/rooms/{roomId}/join_external", participantHandler.JoinRoom).Methods("POST")

	authAPI := api.PathPrefix("/app").Subrouter()
	authAPI.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo))

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

	authAPI.HandleFunc("/rooms/{roomId}/posts", postHandler.CreatePost).Methods("POST")
	authAPI.HandleFunc("/rooms/{roomId}/posts", postHandler.GetPosts).Methods("GET")
	authAPI.HandleFunc("/rooms/{roomId}/posts/{postId}", postHandler.DeletePost).Methods("DELETE")

	authAPI.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
