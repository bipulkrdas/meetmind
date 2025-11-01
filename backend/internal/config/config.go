package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string `env:"PORT" envDefault:"8080"`
	Env                string `env:"ENV" envDefault:"development"`
	DatabaseURL        string `env:"DATABASE_URL,required"`
	JWTSecret          string `env:"JWT_SECRET,required"`
	LiveKitAPIKey      string `env:"LIVEKIT_API_KEY,required"`
	LiveKitAPISecret   string `env:"LIVEKIT_API_SECRET,required"`
	LiveKitURL         string `env:"LIVEKIT_URL,required"`
	EmailProvider      string `env:"EMAIL_PROVIDER" envDefault:"sendgrid"`
	SendGridAPIKey     string `env:"SENDGRID_API_KEY"`
	SendGridFromEmail  string `env:"SENDGRID_FROM_EMAIL"`
	SendGridFromName   string `env:"SENDGRID_FROM_NAME"`
	MailjetAPIKey      string `env:"MAILJET_API_KEY"`
	MailjetSecretKey   string `env:"MAILJET_SECRET_KEY"`
	MailjetFromEmail   string `env:"MAILJET_FROM_EMAIL"`
	MailjetFromName    string `env:"MAILJET_FROM_NAME"`
	FrontendURL        string `env:"FRONTEND_URL,required"`
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`
	FromEmail          string `env:"FROM_EMAIL"`
	FromName           string `env:"FROM_NAME"`
}

func Load() *Config {
	err := godotenv.Load("/Users/bipuldas/livekit/consulting/backend/.env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse config: %+v", err)
	}

	return cfg
}
