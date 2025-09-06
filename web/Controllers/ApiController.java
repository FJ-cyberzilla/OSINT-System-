package com.osint.gateway.controllers;

import com.osint.gateway.models.InvestigationRequest;
import com.osint.gateway.models.OperationResult;
import com.osint.gateway.services.OrchestrationService;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.context.request.async.DeferredResult;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("/api/v1")
public class ApiController {
    
    @Autowired
    private OrchestrationService orchestrationService;
    
    @PostMapping("/investigate")
    public DeferredResult<ResponseEntity<OperationResult>> investigateTarget(
            @Valid @RequestBody InvestigationRequest request) {
        
        DeferredResult<ResponseEntity<OperationResult>> deferredResult = 
            new DeferredResult<>(TimeUnit.SECONDS.toMillis(30));
        
        CompletableFuture<OperationResult> future = orchestrationService.orchestrateInvestigationAsync(
            request.getTarget(),
            request.getPriority(),
            request.getModules()
        );
        
        future.whenComplete((result, ex) -> {
            if (ex != null) {
                deferredResult.setErrorResult(
                    ResponseEntity.internalServerError()
                        .body(OperationResult.error("Investigation failed: " + ex.getMessage()))
                );
            } else {
                deferredResult.setResult(ResponseEntity.ok(result));
            }
        });
        
        return deferredResult;
    }
    
    @GetMapping("/operations/{id}")
    public ResponseEntity<OperationResult> getOperationResult(
            @PathVariable String id) {
        
        OperationResult result = orchestrationService.getOperationResult(id);
        
        if (result.getError() != null && result.getStatus().equals("failed")) {
            return ResponseEntity.status(404).body(result);
        }
        
        return ResponseEntity.ok(result);
    }
    
    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> healthCheck() {
        return ResponseEntity.ok(Map.of(
            "status", "healthy",
            "service", "osint-gateway",
            "timestamp", java.time.Instant.now().toString()
        ));
    }
    
    @DeleteMapping("/cache")
    public ResponseEntity<Map<String, Object>> clearCache() {
        int size = orchestrationService.getCacheSize();
        orchestrationService.clearCache();
        
        return ResponseEntity.ok(Map.of(
            "cleared_entries", size,
            "message", "Cache cleared successfully",
            "timestamp", java.time.Instant.now().toString()
        ));
    }
}
