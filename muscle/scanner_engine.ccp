#include "scanner_engine.hpp"
#include <curl/curl.h>
#include <iostream>
#include <chrono>
#include <iomanip>
#include <sstream>
#include <thread>
#include <random>

// Callback function for writing data
static size_t WriteCallback(void* contents, size_t size, size_t nmemb, std::string* data) {
    data->append((char*)contents, size * nmemb);
    return size * nmemb;
}

ScannerEngine::ScannerEngine() : timeout_ms(5000), is_initialized(false) {
    http_headers = {
        {"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
        {"Accept-Language", "en-US,en;q=0.5"},
        {"Connection", "keep-alive"},
        {"Upgrade-Insecure-Requests", "1"}
    };
}

ScannerEngine::~ScannerEngine() {
    if (is_initialized) {
        curl_global_cleanup();
    }
}

bool ScannerEngine::initialize() {
    if (curl_global_init(CURL_GLOBAL_DEFAULT) != CURLE_OK) {
        std::cerr << "❌ Failed to initialize libcurl" << std::endl;
        return false;
    }
    is_initialized = true;
    std::cout << "✅ ScannerEngine initialized successfully" << std::endl;
    return true;
}

void ScannerEngine::setTimeout(int milliseconds) {
    timeout_ms = milliseconds;
}

std::string ScannerEngine::get_current_timestamp() const {
    auto now = std::chrono::system_clock::now();
    auto in_time_t = std::chrono::system_clock::to_time_t(now);
    std::stringstream ss;
    ss << std::put_time(std::localtime(&in_time_t), "%Y-%m-%d %H:%M:%S");
    return ss.str();
}

bool ScannerEngine::check_status_code(int status_code, const std::vector<int>& success_codes) const {
    if (success_codes.empty()) {
        return status_code == 200;
    }
    
    for (int code : success_codes) {
        if (status_code == code) {
            return true;
        }
    }
    return false;
}

std::string ScannerEngine::make_request(const std::string& url, const std::string& user_agent) {
    if (!is_initialized) {
        return "";
    }

    CURL* curl = curl_easy_init();
    if (!curl) {
        return "";
    }

    std::string response_string;
    std::string header_string;
    long response_code = 0;
    double total_time = 0;

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);
    curl_easy_setopt(curl, CURLOPT_HEADERDATA, &header_string);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, user_agent.empty() ? "MUSCLE-Scanner/1.0" : user_agent.c_str());
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, timeout_ms);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);

    // Add custom headers
    struct curl_slist* headers = nullptr;
    for (const auto& header : http_headers) {
        std::string header_line = header.first + ": " + header.second;
        headers = curl_slist_append(headers, header_line.c_str());
    }
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

    CURLcode res = curl_easy_perform(curl);

    if (res == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response_code);
        curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME, &total_time);
    }

    curl_easy_cleanup(curl);
    curl_slist_free_all(headers);

    return std::to_string(response_code);
}

ScanResult ScannerEngine::check_website(const std::string& target, const json& website_config) {
    ScanResult result;
    
    try {
        std::string url_template = website_config["url"].get<std::string>();
        std::string url = url_template;
        size_t pos = url.find("{}");
        if (pos != std::string::npos) {
            url.replace(pos, 2, target);
        }

        result.target = target;
        result.website_name = website_config["name"].get<std::string>();
        result.url = url;
        result.timestamp = get_current_timestamp();

        // Add random delay to avoid rate limiting
        std::random_device rd;
        std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(100, 800);
        std::this_thread::sleep_for(std::chrono::milliseconds(dis(gen)));

        std::string user_agent = "MUSCLE-Scanner/1.0";
        if (website_config.contains("randomize_user_agent") && website_config["randomize_user_agent"].get<bool>()) {
            std::vector<std::string> user_agents = {
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"
            };
            std::uniform_int_distribution<> ua_dis(0, user_agents.size() - 1);
            user_agent = user_agents[ua_dis(gen)];
        }

        std::string status_str = make_request(url, user_agent);
        int status_code = 0;
        try {
            status_code = std::stoi(status_str);
        } catch (...) {
            status_code = 0;
        }

        result.status_code = status_code;

        std::vector<int> success_codes;
        if (website_config.contains("success_codes")) {
            for (const auto& code : website_config["success_codes"]) {
                success_codes.push_back(code.get<int>());
            }
        } else {
            success_codes = {200};
        }

        result.found = check_status_code(status_code, success_codes);

    } catch (const std::exception& e) {
        std::cerr << "Error checking website: " << e.what() << std::endl;
        result.found = false;
        result.status_code = 0;
    }

    return result;
}

ScanResults ScannerEngine::scan_target(const std::string& target, const std::vector<json>& websites) {
    ScanResults results;
    results.total_checked = websites.size();
    results.found_count = 0;

    auto start_time = std::chrono::high_resolution_clock::now();

    for (const auto& website : websites) {
        try {
            ScanResult result = check_website(target, website);
            if (result.found) {
                results.found_count++;
            }
            results.details.push_back(result);
        } catch (const std::exception& e) {
            std::cerr << "Error processing website: " << e.what() << std::endl;
        }
    }

    auto end_time = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end_time - start_time);

    results.performance_metrics = {
        {"scan_duration_ms", duration.count()},
        {"target", target},
        {"websites_checked", results.total_checked},
        {"websites_found", results.found_count},
        {"timestamp", get_current_timestamp()}
    };

    return results;
}

IntelligenceData ScannerEngine::gather_additional_intel(const std::string& target, const json& scan_data) {
    IntelligenceData intel;
    
    // Placeholder implementation - would integrate with actual APIs
    intel.breach_info = {
        {"target", target},
        {"breaches_found", 0},
        {"message", "Breach data collection would be implemented here"}
    };
    
    intel.whois_info = {
        {"target", target},
        {"domain_info", "Whois data would be retrieved here"},
        {"registration_date", "N/A"},
        {"registrar", "N/A"}
    };
    
    intel.social_analysis = {
        {"target", target},
        {"social_media_presence", "Social media analysis would be performed here"},
        {"influence_score", 0}
    };
    
    intel.threat_data = {
        {"target", target},
        {"threat_level", "unknown"},
        {"malicious_indicators", 0}
    };
    
    return intel;
}

json ScannerEngine::get_performance_metrics() const {
    return {
        {"timeout_ms", timeout_ms},
        {"initialized", is_initialized},
        {"http_headers_count", http_headers.size()}
    };
}
