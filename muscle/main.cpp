#include <zmq.hpp>
#include <string>
#include <iostream>
#include <thread>
#include <atomic>
#include <chrono>
#include <iomanip>
#include <csignal>
#include <vector>
#include <nlohmann/json.hpp>
#include "scanner_engine.hpp"
#include "websites_manager.hpp"
#include "rate-limiter.hpp"

using json = nlohmann::json;

class MuscleEngine {
private:
    zmq::context_t context;
    zmq::socket_t socket;
    ScannerEngine scanner;
    WebsitesManager websites_manager;
    AdaptiveRateLimiter rate_limiter;
    std::atomic<bool> running{false};
    std::atomic<int> consecutive_errors{0};
    
    std::string get_current_timestamp() {
        auto now = std::chrono::system_clock::now();
        auto in_time_t = std::chrono::system_clock::to_time_t(now);
        std::stringstream ss;
        ss << std::put_time(std::localtime(&in_time_t), "%Y-%m-%d %X");
        return ss.str();
    }
    
    // OPSEC: Check for suspicious patterns
    bool is_suspicious_request(const json& message) {
        try {
            std::string target = message["target"].get<std::string>();
            
            // Check for obviously fake or test patterns
            if (target == "test" || target == "admin" || target == "root") {
                return true;
            }
            
            // Check for rapid repeated requests
            static std::string last_target;
            static auto last_time = std::chrono::steady_clock::now();
            auto now = std::chrono::steady_clock::now();
            auto elapsed = std::chrono::duration_cast<std::chrono::milliseconds>(now - last_time).count();
            
            if (target == last_target && elapsed < 500) { // Same target in <500ms
                return true;
            }
            
            last_target = target;
            last_time = now;
            
        } catch (...) {
            // If we can't parse, be cautious
            return true;
        }
        
        return false;
    }
    
public:
    MuscleEngine() : context(1), socket(context, ZMQ_REP), rate_limiter(1500) {
        socket.bind("tcp://*:5556");
        socket.setsockopt(ZMQ_RCVTIMEO, 1000);
        websites_manager.load_websites("websites.json");
        std::cout << "💪 MUSCLE layer initialized with OPSEC features" << std::endl;
        std::cout << "📡 Listening on port 5556" << std::endl;
    }
    
    ~MuscleEngine() {
        stop();
    }
    
    json perform_scan(const std::string& target, const json& pattern_data) {
        // Apply rate limiting with OPSEC
        rate_limiter.wait();
        
        auto websites = websites_manager.get_websites();
        auto results = scanner.scan_target(target, websites);
        
        return {
            {"target", target},
            {"scan_type", "comprehensive"},
            {"websites_checked", websites.size()},
            {"websites_found", results.found_count},
            {"results", results.details},
            {"performance_metrics", scanner.get_performance_metrics()},
            {"timestamp", get_current_timestamp()},
            {"opsec_note", "Rate limited and randomized for operational security"}
        };
    }
    
    json gather_intelligence(const std::string& target, const json& scan_data) {
        // Apply rate limiting with OPSEC
        rate_limiter.wait();
        
        auto intel = scanner.gather_additional_intel(target, scan_data);
        
        return {
            {"target", target},
            {"breach_data", intel.breach_info},
            {"whois_data", intel.whois_info},
            {"social_analysis", intel.social_analysis},
            {"threat_intel", intel.threat_data},
            {"timestamp", get_current_timestamp()},
            {"opsec_note", "Rate limited and randomized for operational security"}
        };
    }
    
    void process_request(const json& message) {
        auto action = message["action"].get<std::string>();
        auto target = message["target"].get<std::string>();
        json response;
        
        try {
            // OPSEC: Check for suspicious requests
            if (is_suspicious_request(message)) {
                std::cout << "🚨 Suspicious request detected: " << target << std::endl;
                rate_limiter.emergency_slowdown();
                response = {{"error", "Request throttled for security reasons"}};
                consecutive_errors++;
            }
            else if (action == "perform_scan") {
                auto pattern_data = message.value("pattern_data", json::object());
                response = perform_scan(target, pattern_data);
                consecutive_errors = 0;
            }
            else if (action == "gather_intelligence") {
                auto scan_data = message.value("scan_data", json::object());
                response = gather_intelligence(target, scan_data);
                consecutive_errors = 0;
            }
            else if (action == "status") {
                response = {
                    {"status", "running"},
                    {"requests_processed", rate_limiter.get_request_count()},
                    {"current_delay_ms", rate_limiter.get_current_delay().count()},
                    {"timestamp", get_current_timestamp()}
                };
            }
            else if (action == "reset_limiter") {
                rate_limiter.reset();
                response = {{"status", "rate_limiter_reset"}};
            }
            else {
                response = {{"error", "Unknown action: " + action}};
                consecutive_errors++;
            }
        } catch (const std::exception& e) {
            response = {{"error", std::string("Processing failed: ") + e.what()}};
            consecutive_errors++;
            
            // Emergency slowdown if too many errors
            if (consecutive_errors > 3) {
                rate_limiter.emergency_slowdown();
            }
        }
        
        zmq::message_t reply(response.dump());
        socket.send(reply, zmq::send_flags::none);
    }
    
    void run() {
        running = true;
        std::cout << "🚀 Muscle engine started with adaptive rate limiting" << std::endl;
        std::cout << "🛡️  OPSEC features enabled: jitter, randomization, pattern detection" << std::endl;
        
        while (running) {
            zmq::message_t request;
            try {
                auto result = socket.recv(request, zmq::recv_flags::none);
                
                if (result) {
                    auto message = json::parse(request.to_string());
                    std::cout << "📨 Received request: " << message["action"] 
                              << " for target: " << message["target"] << std::endl;
                    process_request(message);
                }
            } catch (const zmq::error_t& e) {
                if (e.num() != EAGAIN) {
                    std::cerr << "ZMQ Error: " << e.what() << std::endl;
                }
            } catch (const std::exception& e) {
                std::cerr << "Error processing request: " << e.what() << std::endl;
                json error_response = {{"error", e.what()}};
                zmq::message_t reply(error_response.dump());
                socket.send(reply, zmq::send_flags::none);
            }
        }
    }
    
    void stop() {
        running = false;
    }
};

int main() {
    try {
        MuscleEngine engine;
        
        // Handle graceful shutdown
        std::signal(SIGINT, [](int) {
            std::cout << "\n🛑 Shutting down Muscle engine gracefully..." << std::endl;
            exit(0);
        });
        
        engine.run();
    } catch (const std::exception& e) {
        std::cerr << "❌ Muscle engine failed: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}
