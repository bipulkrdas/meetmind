package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global zerolog logger.
// In "development" environment, it uses a human-friendly console writer.
// In "production" or any other environment, it uses a JSON format suitable for Cloud Run.
func Init(env string) {
	// Default to production settings
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// To better integrate with Google Cloud Logging, which uses "severity"
	zerolog.LevelFieldName = "severity"

	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		// Use a human-friendly console writer for development
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	// For production, the default zerolog behavior (JSON to os.Stderr) is exactly
	// what's needed for Cloud Run, so no special "else" block is required for output.
}
