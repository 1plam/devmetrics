package main

import (
	"devmetrics/pkg/logger"
	"net/http"
)

func main() {
	config := logger.DefaultConfig()
	config.Development = true
	config.Encoding = "console"

	log, err := logger.NewLogger(config)
	if err != nil {
		panic(err)
	}

	// Use logger
	log.Info("Starting application",
		logger.String("app", "devmetrics"),
		logger.String("version", "1.0.0"),
	)

	// Create router and add logging middleware
	router := http.NewServeMux()
	router.HandleFunc("/", handler)

	// Wrap router with logging middleware
	loggedRouter := logger.HTTPMiddleware(log)(router)

	// Start server
	log.Info("Server starting", logger.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", loggedRouter); err != nil {
		log.Fatal("Server failed to start", logger.Error(err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Handler implementation
}
