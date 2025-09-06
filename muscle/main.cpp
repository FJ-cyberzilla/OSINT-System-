#include <zmq.hpp>
#include <string>
#include <iostream>
#include <thread>
#include <atomic>
#include "scanner_engine.hpp"
#include "websites_manager.hpp"

class MuscleEngine {
private:
    zmq::context_t context;
    zmq::socket_t socket;
    ScannerEngine scanner;
    WebsitesManager websites_manager;
    std::atomic<bool> running{false};
    
public:
    MuscleEngine() : context(1), socket(context, ZMQ_REP) {
        socket.bind("tcp://*:5556");
        websites_manager.load_websites("websites.json");
        std::cout << "💪 MUSCLE layer initialized and listening on port 5556" << std::endl;
    }
    
    json perform_scan(const std::string& target, const json& pattern_data) {
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
        // Additional intelligence gathering
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
        
        if (action == "perform_scan") {
            auto pattern_data = message.value("pattern_data", json::object());
            response = perform_scan(target, pattern_data);
        }
        else if (action == "gather_intelligence") {
            auto scan_data = message.value("scan_data", json::object());
            response = gather_intelligence(target, scan_data);
        }
        else {
            response = {{"error", "Unknown action"}};
        }
        
        zmq::message_t reply(response.dump());
        socket.send(reply, zmq::send_flags::none);
    }
    
    void run() {
        running = true;
        while (running) {
            zmq::message_t request;
            auto result = socket.recv(request, zmq::recv_flags::none);
            
            if (result) {
                try {
                    auto message = json::parse(request.to_string());
                    process_request(message);
                } catch (const std::exception& e) {
                    json error_response = {{"error", e.what()}};
                    zmq::message_t reply(error_response.dump());
                    socket.send(reply, zmq::send_flags::none);
                }
            }
        }
    }
};

int main() {
    try {
        MuscleEngine engine;
        engine.run();
    } catch (const std::exception& e) {
        std::cerr << "❌ Muscle engine failed: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}
