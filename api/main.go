package main

import (
	"log"
	"net/http"
	"os"

	"osint-api/handlers"
	"osint-api/handlers/middleware"

	"github.com/gorilla/mux"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	// Initialize ZMQ socket
	zmqSocket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatalf("Failed to create ZMQ socket: %v", err)
	}
	defer zmqSocket.Close()

	orchestraAddr := os.Getenv("ORCHESTRA_ADDR")
	if orchestraAddr == "" {
		orchestraAddr = "tcp://localhost:5558"
	}

	if err := zmqSocket.Connect(orchestraAddr); err != nil {
		log.Fatalf("Failed to connect to orchestra: %v", err)
	}

	// Initialize handlers
	intelHandler := &handlers.IntelHandler{ZmqSocket: zmqSocket}
	healthHandler := &handlers.HealthHandler{ZmqSocket: zmqSocket}
	opsHandler := &handlers.OpsHandler{}

	// Setup router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.AuthMiddleware)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/intel", intelHandler.HandleIntelRequest).Methods("POST")
	api.HandleFunc("/intel/batch", intelHandler.HandleBatchIntelRequest).Methods("POST")
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	api.HandleFunc("/ready", healthHandler.ReadyCheck).Methods("GET")
	api.HandleFunc("/stats", healthHandler.StatsHandler).Methods("GET")
	api.HandleFunc("/operations", opsHandler.ListOperations).Methods("GET")
	api.HandleFunc("/operations/status", opsHandler.GetOperationStatus).Methods("GET")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌐 OSINT API server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
