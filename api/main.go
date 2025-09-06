package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "github.com/gorilla/mux"
    zmq "github.com/pebbe/zmq4"
)

type IntelRequest struct {
    Target    string                 `json:"target"`
    ScanData  map[string]interface{} `json:"scan_data"`
    Timestamp time.Time              `json:"timestamp"`
    OperationID string               `json:"operation_id"`
}

type APIHandler struct {
    zmqSocket *zmq.Socket
    client    *http.Client
}

func (h *APIHandler) HandleIntelRequest(w http.ResponseWriter, r *http.Request) {
    var req IntelRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
        return
    }

    if req.Target == "" {
        http.Error(w, `{"error": "Target is required"}`, http.StatusBadRequest)
        return
    }

    // Set operation ID if not provided
    if req.OperationID == "" {
        req.OperationID = fmt.Sprintf("op_%d", time.Now().Unix())
    }
    req.Timestamp = time.Now()

    // Forward to Orchestra layer
    message := map[string]interface{}{
        "action": "investigate",
        "target": req.Target,
        "scan_data": req.ScanData,
        "operation_id": req.OperationID,
        "timestamp": req.Timestamp,
    }

    if _, err := h.zmqSocket.SendMessage(message); err != nil {
        http.Error(w, `{"error": "Failed to communicate with orchestra"}`, http.StatusInternalServerError)
        return
    }

    reply, err := h.zmqSocket.RecvMessage(0)
    if err != nil {
        http.Error(w, `{"error": "Failed to receive response from orchestra"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Operation-ID", req.OperationID)
    w.WriteHeader(http.StatusOK)
    
    // Send the raw reply from orchestra
    if len(reply) > 0 {
        w.Write([]byte(reply[0]))
    } else {
        w.Write([]byte(`{"error": "Empty response from orchestra"}`))
    }
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().UTC(),
        "version":   "1.0.0",
        "service":   "osint-api",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Connect to Orchestra via ZeroMQ
    zmqSocket, err := zmq.NewSocket(zmq.REQ)
    if err != nil {
        log.Fatalf("Failed to create ZMQ socket: %v", err)
    }
    defer zmqSocket.Close()

    // Connect to orchestra service (use environment variable for flexibility)
    orchestraAddr := os.Getenv("ORCHESTRA_ADDR")
    if orchestraAddr == "" {
        orchestraAddr = "tcp://localhost:5558"
    }

    if err := zmqSocket.Connect(orchestraAddr); err != nil {
        log.Fatalf("Failed to connect to orchestra at %s: %v", orchestraAddr, err)
    }

    // Set timeout for ZMQ operations
    zmqSocket.SetRcvtimeo(30 * time.Second)
    zmqSocket.SetSndtimeo(10 * time.Second)

    // Setup HTTP server
    router := mux.NewRouter()
    handler := &APIHandler{zmqSocket: zmqSocket}

    // API routes
    router.HandleFunc("/api/v1/intel", handler.HandleIntelRequest).Methods("POST")
    router.HandleFunc("/api/v1/health", healthCheck).Methods("GET")
    router.HandleFunc("/api/v1/ready", healthCheck).Methods("GET")

    // Get port from environment or default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🌐 OSINT API layer running on :%s", port)
    log.Printf("🔗 Connected to orchestra at: %s", orchestraAddr)
    
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
