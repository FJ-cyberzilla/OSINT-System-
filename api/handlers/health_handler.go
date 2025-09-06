package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type HealthHandler struct {
	ZmqSocket *zmq.Socket
}

type HealthResponse struct {
	Status     string                 `json:"status"`
	Timestamp  time.Time              `json:"timestamp"`
	Version    string                 `json:"version"`
	System     SystemInfo             `json:"system"`
	Components map[string]string      `json:"components"`
	Uptime     string                 `json:"uptime"`
}

type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	NumCPU       int    `json:"num_cpu"`
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check ZMQ connection health
	zmqStatus := "healthy"
	if h.ZmqSocket == nil {
		zmqStatus = "disconnected"
	} else {
		// Test ZMQ connection by sending a ping
		testMsg := map[string]interface{}{
			"action":  "ping",
			"timeout": 1,
		}
		if _, err := h.ZmqSocket.SendMessage(testMsg); err != nil {
			zmqStatus = "unhealthy"
		}
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			NumCPU:       runtime.NumCPU(),
		},
		Components: map[string]string{
			"zmq_connection": zmqStatus,
			"http_server":    "healthy",
			"memory":         "healthy",
		},
		Uptime: getUptime(),
	}

	// If ZMQ is unhealthy, overall status is degraded
	if zmqStatus != "healthy" {
		response.Status = "degraded"
		response.Components["zmq_connection"] = zmqStatus
	}

	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
		"services": map[string]bool{
			"zmq_connected": h.ZmqSocket != nil,
			"http_listening": true,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"memory":    getMemoryStats(),
		"goroutines": runtime.NumGoroutine(),
		"system": map[string]interface{}{
			"cpu_cores": runtime.NumCPU(),
			"go_version": runtime.Version(),
		},
	}

	json.NewEncoder(w).Encode(stats)
}

func getUptime() string {
	// This would be implemented to track actual process uptime
	return "0h5m" // Example
}

func getMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return map[string]interface{}{
		"allocated":      m.Alloc,
		"total_alloc":    m.TotalAlloc,
		"system":         m.Sys,
		"garbage_collector": m.NumGC,
	}
}
