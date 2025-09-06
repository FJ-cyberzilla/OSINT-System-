package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Operation represents an OSINT investigation operation
type Operation struct {
	ID           string                 `json:"id"`
	Target       string                 `json:"target"`
	Status       string                 `json:"status"` // pending, processing, completed, failed, cancelled
	Priority     string                 `json:"priority"` // low, medium, high, critical
	Progress     float64                `json:"progress"` // 0-100
	CreatedAt    time.Time              `json:"created_at"`
	StartedAt    *time.Time             `json:"started_at,omitempty"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	Results      map[string]interface{} `json:"results,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Duration     string                 `json:"duration,omitempty"`
	Resources    []string               `json:"resources,omitempty"` // Scrapy, SpiderFoot, etc.
	RiskScore    float64                `json:"risk_score,omitempty"`
	Findings     int                    `json:"findings_count,omitempty"`
}

// OpsHandler manages OSINT operations
type OpsHandler struct {
	operations map[string]*Operation
	mu         sync.RWMutex
}

// NewOpsHandler creates a new operations handler
func NewOpsHandler() *OpsHandler {
	return &OpsHandler{
		operations: make(map[string]*Operation),
	}
}

// CreateOperation creates a new OSINT operation
func (h *OpsHandler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		Target   string `json:"target"`
		Priority string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Target == "" {
		h.sendError(w, "Target is required", http.StatusBadRequest)
		return
	}

	if request.Priority == "" {
		request.Priority = "medium"
	}

	operationID := generateOperationID()
	now := time.Now()

	operation := &Operation{
		ID:        operationID,
		Target:    request.Target,
		Status:    "pending",
		Priority:  request.Priority,
		Progress:  0,
		CreatedAt: now,
		Resources: []string{"Scrapy", "SpiderFoot", "AI Analysis"},
	}

	h.mu.Lock()
	h.operations[operationID] = operation
	h.mu.Unlock()

	// Simulate async processing (in real implementation, this would be a background worker)
	go h.processOperation(operation)

	response := map[string]interface{}{
		"operation_id": operationID,
		"status":       "created",
		"message":      "Operation queued for processing",
		"created_at":   now,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetOperationStatus returns the status of a specific operation
func (h *OpsHandler) GetOperationStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	operationID := r.URL.Query().Get("id")
	if operationID == "" {
		h.sendError(w, "Operation ID is required", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	operation, exists := h.operations[operationID]
	h.mu.RUnlock()

	if !exists {
		h.sendError(w, "Operation not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(operation)
}

// ListOperations returns a list of all operations with optional filtering
func (h *OpsHandler) ListOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters for filtering
	statusFilter := r.URL.Query().Get("status")
	priorityFilter := r.URL.Query().Get("priority")
	limit := 50 // Default limit

	// Parse limit parameter
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if n, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil || n != 1 {
			limit = 50
		}
		if limit > 1000 {
			limit = 1000
		}
	}

	h.mu.RLock()
	operations := make([]*Operation, 0, len(h.operations))
	
	for _, op := range h.operations {
		// Apply filters
		if statusFilter != "" && op.Status != statusFilter {
			continue
		}
		if priorityFilter != "" && op.Priority != priorityFilter {
			continue
		}
		operations = append(operations, op)
	}
	h.mu.RUnlock()

	// Sort by creation time (newest first)
	for i := 0; i < len(operations)-1; i++ {
		for j := i + 1; j < len(operations); j++ {
			if operations[i].CreatedAt.Before(operations[j].CreatedAt) {
				operations[i], operations[j] = operations[j], operations[i]
			}
		}
	}

	// Apply limit
	if len(operations) > limit {
		operations = operations[:limit]
	}

	response := map[string]interface{}{
		"operations": operations,
		"total":      len(operations),
		"limit":      limit,
		"filters": map[string]string{
			"status":   statusFilter,
			"priority": priorityFilter,
		},
		"timestamp": time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// CancelOperation cancels a running operation
func (h *OpsHandler) CancelOperation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	operationID := r.URL.Query().Get("id")
	if operationID == "" {
		h.sendError(w, "Operation ID is required", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	operation, exists := h.operations[operationID]
	if exists && (operation.Status == "pending" || operation.Status == "processing") {
		operation.Status = "cancelled"
		operation.CompletedAt = &[]time.Time{time.Now()}[0]
		operation.Progress = 0
		operation.Error = "Operation cancelled by user"
	}
	h.mu.Unlock()

	if !exists {
		h.sendError(w, "Operation not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"operation_id": operationID,
		"status":       "cancelled",
		"message":      "Operation cancelled successfully",
		"cancelled_at": time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// GetOperationsStats returns statistics about operations
func (h *OpsHandler) GetOperationsStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := map[string]interface{}{
		"total_operations":    len(h.operations),
		"pending_operations":  0,
		"processing_operations": 0,
		"completed_operations": 0,
		"failed_operations":   0,
		"cancelled_operations": 0,
		"average_duration":    "0s",
		"success_rate":        0.0,
	}

	var totalDuration time.Duration
	var completedCount int

	for _, op := range h.operations {
		switch op.Status {
		case "pending":
			stats["pending_operations"] = stats["pending_operations"].(int) + 1
		case "processing":
			stats["processing_operations"] = stats["processing_operations"].(int) + 1
		case "completed":
			stats["completed_operations"] = stats["completed_operations"].(int) + 1
			completedCount++
			if op.StartedAt != nil && op.CompletedAt != nil {
				totalDuration += op.CompletedAt.Sub(*op.StartedAt)
			}
		case "failed":
			stats["failed_operations"] = stats["failed_operations"].(int) + 1
		case "cancelled":
			stats["cancelled_operations"] = stats["cancelled_operations"].(int) + 1
		}
	}

	// Calculate success rate
	totalCompleted := stats["completed_operations"].(int) + stats["failed_operations"].(int) + stats["cancelled_operations"].(int)
	if totalCompleted > 0 {
		stats["success_rate"] = float64(stats["completed_operations"].(int)) / float64(totalCompleted) * 100
	}

	// Calculate average duration
	if completedCount > 0 {
		avgDuration := totalDuration / time.Duration(completedCount)
		stats["average_duration"] = avgDuration.String()
	}

	stats["timestamp"] = time.Now()
	json.NewEncoder(w).Encode(stats)
}

// CleanupOperations removes old operations
func (h *OpsHandler) CleanupOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	maxAgeStr := r.URL.Query().Get("max_age")
	if maxAgeStr == "" {
		maxAgeStr = "720h" // 30 days default
	}

	maxAge, err := time.ParseDuration(maxAgeStr)
	if err != nil {
		h.sendError(w, "Invalid duration format", http.StatusBadRequest)
		return
	}

	cutoff := time.Now().Add(-maxAge)
	deletedCount := 0

	h.mu.Lock()
	for id, op := range h.operations {
		if op.CreatedAt.Before(cutoff) && (op.Status == "completed" || op.Status == "failed" || op.Status == "cancelled") {
			delete(h.operations, id)
			deletedCount++
		}
	}
	h.mu.Unlock()

	response := map[string]interface{}{
		"deleted_count": deletedCount,
		"max_age":       maxAge.String(),
		"cutoff_time":   cutoff,
		"remaining_operations": len(h.operations),
		"timestamp":     time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// processOperation simulates processing an OSINT operation
func (h *OpsHandler) processOperation(operation *Operation) {
	h.mu.Lock()
	operation.Status = "processing"
	startTime := time.Now()
	operation.StartedAt = &startTime
	h.mu.Unlock()

	// Simulate processing stages
	stages := []struct {
		name     string
		duration time.Duration
		progress float64
	}{
		{"Initializing", 2 * time.Second, 10},
		{"Scrapy Crawling", 8 * time.Second, 30},
		{"SpiderFoot Analysis", 6 * time.Second, 60},
		{"AI Pattern Recognition", 4 * time.Second, 80},
		{"Final Correlation", 2 * time.Second, 100},
	}

	for _, stage := range stages {
		time.Sleep(stage.duration)

		h.mu.Lock()
		if operation.Status == "cancelled" {
			h.mu.Unlock()
			return
		}
		operation.Progress = stage.progress
		h.mu.Unlock()
	}

	// Final results
	h.mu.Lock()
	operation.Status = "completed"
	completeTime := time.Now()
	operation.CompletedAt = &completeTime
	operation.Duration = completeTime.Sub(startTime).String()
	operation.Results = map[string]interface{}{
		"platforms_checked": 15,
		"findings_count":    23,
		"risk_score":        0.72,
		"pattern_type":      "professional_account",
		"confidence":        0.88,
	}
	operation.Findings = 23
	operation.RiskScore = 0.72
	h.mu.Unlock()
}

// sendError sends a standardized error response
func (h *OpsHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]interface{}{
		"error":       message,
		"status":      "error",
		"status_code": statusCode,
		"timestamp":   time.Now(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// generateOperationID generates a unique operation ID
func generateOperationID() string {
	return fmt.Sprintf("op_%d_%s", time.Now().Unix(), randomString(6))
}

// randomString generates a random string for operation IDs
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
