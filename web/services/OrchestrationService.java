package com.osint.gateway.services;

import com.osint.gateway.models.OperationResult;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;

import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class OrchestrationService {
    
    @Value("${orchestrator.url:http://localhost:5558}")
    private String orchestratorUrl;
    
    @Value("${orchestrator.timeout:30000}")
    private int timeout;
    
    private final RestTemplate restTemplate;
    private final Map<String, OperationResult> operationCache;
    
    public OrchestrationService() {
        this.restTemplate = new RestTemplate();
        this.operationCache = new ConcurrentHashMap<>();
    }
    
    @Retryable(
        value = {ResourceAccessException.class, HttpServerErrorException.class},
        maxAttempts = 3,
        backoff = @Backoff(delay = 1000, multiplier = 2)
    )
    public OperationResult orchestrateInvestigation(String target, Integer priority, java.util.List<String> modules) {
        String url = UriComponentsBuilder.fromHttpUrl(orchestratorUrl)
                .path("/orchestrate")
                .build()
                .toUriString();
        
        // Prepare request body
        Map<String, Object> requestBody = Map.of(
            "target", target,
            "priority", priority != null ? priority : 1,
            "modules", modules != null ? modules : java.util.List.of("all"),
            "timestamp", java.time.Instant.now().toString()
        );
        
        // Set headers
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.set("User-Agent", "OSINT-Gateway/1.0.0");
        
        HttpEntity<Map<String, Object>> entity = new HttpEntity<>(requestBody, headers);
        
        try {
            ResponseEntity<Map> response = restTemplate.exchange(
                url,
                HttpMethod.POST,
                entity,
                Map.class
            );
            
            if (response.getStatusCode() == HttpStatus.OK && response.getBody() != null) {
                return mapToOperationResult(response.getBody());
            } else {
                return OperationResult.error("Orchestrator returned unexpected response: " + response.getStatusCode());
            }
            
        } catch (HttpClientErrorException e) {
            return OperationResult.error("Client error when calling orchestrator: " + e.getMessage());
        } catch (HttpServerErrorException e) {
            return OperationResult.error("Server error from orchestrator: " + e.getMessage());
        } catch (ResourceAccessException e) {
            return OperationResult.error("Cannot connect to orchestrator at: " + orchestratorUrl);
        } catch (Exception e) {
            return OperationResult.error("Unexpected error: " + e.getMessage());
        }
    }
    
    @Retryable(
        value = {ResourceAccessException.class},
        maxAttempts = 2,
        backoff = @Backoff(delay = 500)
    )
    public OperationResult getOperationResult(String operationId) {
        // Check cache first
        OperationResult cachedResult = operationCache.get(operationId);
        if (cachedResult != null && "completed".equals(cachedResult.getStatus())) {
            return cachedResult;
        }
        
        String url = UriComponentsBuilder.fromHttpUrl(orchestratorUrl)
                .path("/operations/" + operationId)
                .build()
                .toUriString();
        
        try {
            ResponseEntity<Map> response = restTemplate.getForEntity(url, Map.class);
            
            if (response.getStatusCode() == HttpStatus.OK && response.getBody() != null) {
                OperationResult result = mapToOperationResult(response.getBody());
                
                // Cache completed operations
                if ("completed".equals(result.getStatus()) || "failed".equals(result.getStatus())) {
                    operationCache.put(operationId, result);
                }
                
                return result;
            } else {
                return OperationResult.error("Orchestrator returned unexpected response for operation: " + operationId);
            }
            
        } catch (HttpClientErrorException.NotFound e) {
            return OperationResult.error("Operation not found: " + operationId);
        } catch (Exception e) {
            return OperationResult.error("Error fetching operation result: " + e.getMessage());
        }
    }
    
    public CompletableFuture<OperationResult> orchestrateInvestigationAsync(
            String target, Integer priority, java.util.List<String> modules) {
        return CompletableFuture.supplyAsync(() -> 
            orchestrateInvestigation(target, priority, modules)
        );
    }
    
    public void clearCache() {
        operationCache.clear();
    }
    
    public int getCacheSize() {
        return operationCache.size();
    }
    
    private OperationResult mapToOperationResult(Map<String, Object> response) {
        OperationResult result = new OperationResult();
        
        result.setOperationId((String) response.get("operation_id"));
        result.setTarget((String) response.get("target"));
        result.setStatus((String) response.get("status"));
        
        if (response.get("progress") != null) {
            result.setProgress(Double.valueOf(response.get("progress").toString()));
        }
        
        if (response.get("results") instanceof Map) {
            result.setResults((Map<String, Object>) response.get("results"));
        }
        
        if (response.get("error") != null) {
            result.setError(response.get("error").toString());
        }
        
        return result;
    }
}
