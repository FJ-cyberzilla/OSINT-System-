#include <zmq.hpp>
#include <string>
#include <iostream>
#include <thread>
#include <atomic>
#include <chrono>
#include <iomanip>
#include <csignal>
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
    RateLimiter rate_limiter;
    std::atomic<bool> running{false};
    
    std::string get_current_timestamp() {
        auto now = std::chrono::system_clock::now();
        auto in_time_t = std::chrono::system_clock::to_time_t(now);
        std::stringstream ss;
        ss << std::put_time(std::localtime(&in_time_t), "%Y-%m-%d %X");
        return ss.str();
    }
    
public:
    MuscleEngine() : context(1), socket(context, ZMQ_REP), rate_limiter(10, 1000) {
        socket.bind("tcp://*:5556");
        socket.setsockopt(ZMQ_RCVTIMEO, 1000);
        websites_manager.load_websites("websites.json");
        std::cout << "💪 MUSCLE layer initialized and listening on port 5556" << std::endl;
    }
    
    ~MuscleEngine() {
        stop();
    }
    
    json perform_scan(const std::string& target, const json& pattern_data) {
        if (!rate_limiter.allow_request()) {
            throw std::runtime_error("Rate limit exceeded. Please try again later.");
        }
        
        auto websites = websites_manager.get_websites();
        auto results = scanner.scan_target(target, websites);
        
        return {
            {"target", target},
            {"scan_type", "comprehensive"},
            {"websites_checked", websites.size()},
            {"websites_found", results.found_count},
            {"results", results.details},
            {"performance_metrics", scanner.get_performance_metrics()},
            {"timestamp", get_current_timestamp()}
        };
    }
    
    json gather_intelligence(const std::string& target, const json& scan_data) {
        if (!rate_limiter.allow_request()) {
            throw std::runtime_error("Rate limit exceeded. Please try again later.");
        }
        
        auto intel = scanner.gather_additional_intel(target, scan_data);
        
        return {
            {"target", target},
            {"breach_data", intel.breach_info},
            {"whois_data", intel.whois_info},
            {"social_analysis", intel.social_analysis},
            {"threat_intel", intel.threat_data},
            {"timestamp", get_current_timestamp()}
        };
    }
    
    void process_request(const json& message) {
        auto action = message["action"].get<std::string>();
        auto target = message["target"].get<std::string>();
        json response;
        
        try {
            if (action == "perform_scan") {
                auto pattern_data = message.value("pattern_data", json::object());
                response = perform_scan(target, pattern_data);
            }
            else if (action == "gather_intelligence") {
                auto scan_data = message.value("scan_data", json::object());
                response = gather_intelligence(target, scan_data);
            }
            else if (action == "status") {
                response = {
                    {"status", "running"},
                    {"requests_processed", rate_limiter.get_request_count()},
                    {"timestamp", get_current_timestamp()}
                };
            }
            else {
                response = {{"error", "Unknown action: " + action}};
            }
        } catch (const std::exception& e) {
            response = {{"error", std::string("Processing failed: ") + e.what()}};
        }
        
        zmq::message_t reply(response.dump());
        socket.send(reply, zmq::send_flags::none);
    }
    
    void run() {
        running = true;
        std::cout << "🚀 Muscle engine started. Waiting for requests..." << std::endl;
        
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
            std::cout << "\n🛑 Shutting down Muscle engine..." << std::endl;
            exit(0);
        });
        
        engine.run();
    } catch (const std::exception& e) {
        std::cerr << "❌ Muscle engine failed: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}
