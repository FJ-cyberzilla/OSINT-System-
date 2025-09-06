package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"

	zmq "github.com/pebbe/zmq4"
)

type IntelHandler struct {
	ZmqSocket *zmq.Socket
}

type IntelRequest struct {
	Target      string                 `json:"target"`
	ScanData    map[string]interface{} `json:"scan_data"`
	OperationID string                 `json:"operation_id"`
	Priority    string                 `json:"priority"` // low, medium, high
}

type IntelResponse struct {
	OperationID  string                 `json:"operation_id"`
	Target       string                 `json:"target"`
	Status       string                 `json:"status"`
	Results      map[string]interface{} `json:"results"`
	Timestamps   Timestamps             `json:"timestamps"`
	RiskAssessment *RiskAssessment      `json:"risk_assessment,omitempty"`
}

type Timestamps struct {
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`
	Duration string    `json:"duration"`
}

type RiskAssessment struct {
	Score         float64  `json:"score"`
	Level         string   `json:"level"`
	Factors       []string `json:"factors"`
	Confidence    float64  `json:"confidence"`
	Recommendations []string `json:"recommendations"`
}

func (h *IntelHandler) HandleIntelRequest(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	var req IntelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Target == "" {
		h.sendError(w, "Target is required", http.StatusBadRequest)
		return
	}

	// Set operation ID if not provided
	if req.OperationID == "" {
		req.OperationID = generateOperationID()
	}

	// Set default priority
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Prepare message for Orchestra
	message := map[string]interface{}{
		"action":       "investigate",
		"target":       req.Target,
		"scan_data":    req.ScanData,
		"operation_id": req.OperationID,
		"priority":     req.Priority,
		"timestamp":    time.Now(),
	}

	// Send to Orchestra
	if _, err := h.ZmqSocket.SendMessage(message); err != nil {
		h.sendError(w, "Failed to communicate with orchestra", http.StatusInternalServerError)
		return
	}

	// Receive response
	reply, err := h.ZmqSocket.RecvMessage(0)
	if err != nil {
		h.sendError(w, "Failed to receive response from orchestra", http.StatusInternalServerError)
		return
	}

	// Process response
	if len(reply) == 0 {
		h.sendError(w, "Empty response from orchestra", http.StatusInternalServerError)
		return
	}

	// Forward the orchestra response
	w.Header().Set("X-Operation-ID", req.OperationID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reply[0]))
}

func (h *IntelHandler) HandleBatchIntelRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requests []IntelRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		h.sendError(w, "Invalid JSON format for batch request", http.StatusBadRequest)
		return
	}

	if len(requests) == 0 {
		h.sendError(w, "Empty batch request", http.StatusBadRequest)
		return
	}

	if len(requests) > 100 {
		h.sendError(w, "Batch size too large (max 100)", http.StatusBadRequest)
		return
	}

	results := make([]map[string]interface{}, len(requests))
	for i, req := range requests {
		if req.OperationID == "" {
			req.OperationID = generateOperationID()
		}

		message := map[string]interface{}{
			"action":       "investigate",
			"target":       req.Target,
			"operation_id": req.OperationID,
			"priority":     req.Priority,
			"timestamp":    time.Now(),
			"batch_index":  i,
		}

		if _, err := h.ZmqSocket.SendMessage(message); err != nil {
			results[i] = map[string]interface{}{
				"operation_id": req.OperationID,
				"status":       "error",
				"error":        "Failed to send to orchestra",
			}
			continue
		}

		reply, err := h.ZmqSocket.RecvMessage(0)
		if err != nil {
			results[i] = map[string]interface{}{
				"operation_id": req.OperationID,
				"status":       "error",
				"error":        "Failed to receive response",
			}
			continue
		}

		if len(reply) > 0 {
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(reply[0]), &result); err != nil {
				results[i] = map[string]interface{}{
					"operation_id": req.OperationID,
					"status":       "error",
					"error":        "Invalid response format",
				}
			} else {
				results[i] = result
			}
		}
	}

	response := map[string]interface{}{
		"batch_id":    generateOperationID(),
		"total":       len(requests),
		"successful":  countSuccessful(results),
		"failed":      countFailed(results),
		"operations":  results,
		"timestamp":   time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

func (h *IntelHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]interface{}{
		"error":       message,
		"status":      "error",
		"timestamp":   time.Now(),
		"status_code": statusCode,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func generateOperationID() string {
	return fmt.Sprintf("op_%d_%s", time.Now().Unix(), randomString(6))
}

func randomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "default"
	}
	for i := range bytes {
		bytes[i] = byte(int(bytes[i])%26 + 97) // a-z
	}
	return string(bytes)
}

func countSuccessful(results []map[string]interface{}) int {
	count := 0
	for _, result := range results {
		if status, ok := result["status"].(string); ok && status != "error" {
			count++
		}
	}
	return count
}

func countFailed(results []map[string]interface{}) int {
	count := 0
	for _, result := range results {
		if status, ok := result["status"].(string); ok && status == "error" {
			count++
		}
	}
	return count
}
