#ifndef SCANNER_ENGINE_HPP
#define SCANNER_ENGINE_HPP

#include <string>
#include <vector>
#include <map>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

struct ScanResult {
    std::string target;
    std::string website_name;
    std::string url;
    bool found;
    int status_code;
    std::string response_time;
    std::string timestamp;
};

struct ScanResults {
    int found_count;
    int total_checked;
    std::vector<ScanResult> details;
    json performance_metrics;
};

struct IntelligenceData {
    json breach_info;
    json whois_info;
    json social_analysis;
    json threat_data;
};

class ScannerEngine {
public:
    ScannerEngine();
    ~ScannerEngine();
    
    bool initialize();
    ScanResults scan_target(const std::string& target, const std::vector<json>& websites);
    IntelligenceData gather_additional_intel(const std::string& target, const json& scan_data);
    void setTimeout(int milliseconds);
    json get_performance_metrics() const;
    
private:
    int timeout_ms;
    bool is_initialized;
    std::map<std::string, std::string> http_headers;
    
    ScanResult check_website(const std::string& target, const json& website_config);
    std::string make_request(const std::string& url, const std::string& user_agent = "");
    bool check_status_code(int status_code, const std::vector<int>& success_codes) const;
    std::string get_current_timestamp() const;
};

#endif // SCANNER_ENGINE_HPP
