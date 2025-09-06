package com.osint.gateway.models;

import com.fasterxml.jackson.annotation.JsonInclude;
import java.time.LocalDateTime;
import java.util.Map;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class OperationResult {
    
    private String operationId;
    private String target;
    private String status; // pending, processing, completed, failed
    private Double progress; // 0-100
    private LocalDateTime createdAt;
    private LocalDateTime startedAt;
    private LocalDateTime completedAt;
    private Map<String, Object> results;
    private String error;
    private Map<String, Object> metadata;
    
    // Static factory methods
    public static OperationResult success(String operationId, String target, Map<String, Object> results) {
        OperationResult result = new OperationResult();
        result.setOperationId(operationId);
        result.setTarget(target);
        result.setStatus("completed");
        result.setProgress(100.0);
        result.setResults(results);
        result.setCompletedAt(LocalDateTime.now());
        return result;
    }
    
    public static OperationResult error(String errorMessage) {
        OperationResult result = new OperationResult();
        result.setStatus("failed");
        result.setError(errorMessage);
        result.setCompletedAt(LocalDateTime.now());
        return result;
    }
    
    public static OperationResult processing(String operationId, String target) {
        OperationResult result = new OperationResult();
        result.setOperationId(operationId);
        result.setTarget(target);
        result.setStatus("processing");
        result.setProgress(0.0);
        result.setStartedAt(LocalDateTime.now());
        return result;
    }
    
    // Getters and Setters
    public String getOperationId() { return operationId; }
    public void setOperationId(String operationId) { this.operationId = operationId; }
    
    public String getTarget() { return target; }
    public void setTarget(String target) { this.target = target; }
    
    public String getStatus() { return status; }
    public void setStatus(String status) { this.status = status; }
    
    public Double getProgress() { return progress; }
    public void setProgress(Double progress) { this.progress = progress; }
    
    public LocalDateTime getCreatedAt() { return createdAt; }
    public void setCreatedAt(LocalDateTime createdAt) { this.createdAt = createdAt; }
    
    public LocalDateTime getStartedAt() { return startedAt; }
    public void setStartedAt(LocalDateTime startedAt) { this.startedAt = startedAt; }
    
    public LocalDateTime getCompletedAt() { return completedAt; }
    public void setCompletedAt(LocalDateTime completedAt) { this.completedAt = completedAt; }
    
    public Map<String, Object> getResults() { return results; }
    public void setResults(Map<String, Object> results) { this.results = results; }
    
    public String getError() { return error; }
    public void setError(String error) { this.error = error; }
    
    public Map<String, Object> getMetadata() { return metadata; }
    public void setMetadata(Map<String, Object> metadata) { this.metadata = metadata; }
}
