package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    zmq "github.com/pebbe/zmq4"
)

type IntelRequest struct {
    Target    string                 `json:"target"`
    ScanData  map[string]interface{} `json:"scan_data"`
    Timestamp time.Time              `json:"timestamp"`
}

type APIHandler struct {
    zmqSocket *zmq.Socket
    client    *http.Client
}

func (h *APIHandler) HandleIntelRequest(w http.ResponseWriter, r *http.Request) {
    var req IntelRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Forward to Orchestra layer
    message := map[string]interface{}{
        "action": "investigate",
        "target": req.Target,
        "scan_data": req.ScanData,
    }
    
    h.zmqSocket.SendMessage(message)
    reply, _ := h.zmqSocket.RecvMessage(0)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reply)
}

func main() {
    // Connect to Orchestra
    zmqSocket, _ := zmq.NewSocket(zmq.REQ)
    defer zmqSocket.Close()
    zmqSocket.Connect("tcp://localhost:5558")
    
    // Setup HTTP server
    router := mux.NewRouter()
    handler := &APIHandler{zmqSocket: zmqSocket}
    
    router.HandleFunc("/api/v1/intel", handler.HandleIntelRequest).Methods("POST")
    router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }).Methods("GET")
    
    log.Println("🌐 API layer running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
