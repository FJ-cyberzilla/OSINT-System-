package com.osint.gateway.models;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;
import java.util.List;

public class InvestigationRequest {
    
    @NotBlank(message = "Target is required")
    @Size(min = 2, max = 100, message = "Target must be between 2 and 100 characters")
    private String target;
    
    @NotNull(message = "Priority is required")
    private Integer priority;
    
    private List<String> modules;
    
    private String operationId;
    
    private String callbackUrl;
    
    // Constructors
    public InvestigationRequest() {}
    
    public InvestigationRequest(String target, Integer priority, List<String> modules) {
        this.target = target;
        this.priority = priority;
        this.modules = modules;
    }
    
    // Getters and Setters
    public String getTarget() { return target; }
    public void setTarget(String target) { this.target = target; }
    
    public Integer getPriority() { return priority; }
    public void setPriority(Integer priority) { this.priority = priority; }
    
    public List<String> getModules() { return modules; }
    public void setModules(List<String> modules) { this.modules = modules; }
    
    public String getOperationId() { return operationId; }
    public void setOperationId(String operationId) { this.operationId = operationId; }
    
    public String getCallbackUrl() { return callbackUrl; }
    public void setCallbackUrl(String callbackUrl) { this.callbackUrl = callbackUrl; }
}
