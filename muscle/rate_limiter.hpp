// Adaptive Rate Limiting for OPSEC
#ifndef RATE_LIMITER_HPP
#define RATE_LIMITER_HPP

#include <chrono>
#include <random>
#include <atomic>
#include <mutex>
#include <thread>
#include <iostream>

class AdaptiveRateLimiter {
private:
    std::chrono::milliseconds base_delay;
    std::chrono::milliseconds current_delay;
    std::atomic<int> request_count;
    std::mutex mutex;
    std::default_random_engine generator;
    std::uniform_real_distribution<double> distribution;
    
    // OPSEC: Random user agents to avoid fingerprinting
    std::vector<std::string> user_agents = {
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15",
        "Mozilla/5.0 (Android 10; Mobile; rv:91.0) Gecko/91.0 Firefox/91.0"
    };
    
    std::uniform_int_distribution<size_t> ua_distribution;
    
public:
    AdaptiveRateLimiter(int base_delay_ms = 1000) 
        : base_delay(base_delay_ms),
          current_delay(base_delay_ms),
          request_count(0),
          distribution(0.8, 1.2),
          ua_distribution(0, user_agents.size() - 1) {
        // Seed with current time for proper randomness
        generator.seed(std::chrono::system_clock::now().time_since_epoch().count());
    }
    
    void wait() {
        std::lock_guard<std::mutex> lock(mutex);
        
        // Add jitter to avoid pattern detection
        double jitter = distribution(generator);
        auto actual_delay = std::chrono::milliseconds(
            static_cast<int>(current_delay.count() * jitter)
        );
        
        std::this_thread::sleep_for(actual_delay);
        
        // Adaptive rate limiting - increase delay after each request
        request_count++;
        if (request_count % 5 == 0) { // Adjust more frequently for OPSEC
            current_delay = std::chrono::milliseconds(
                std::min(8000, static_cast<int>(current_delay.count() * 1.15)) // More conservative
            );
            std::cout << "🔧 Rate limit adjusted to: " << current_delay.count() << "ms" << std::endl;
        }
    }
    
    // OPSEC: Get random user agent to avoid detection
    std::string get_random_user_agent() {
        std::lock_guard<std::mutex> lock(mutex);
        return user_agents[ua_distribution(generator)];
    }
    
    // OPSEC: Get current delay with jitter for external use
    std::chrono::milliseconds get_jittered_delay() {
        std::lock_guard<std::mutex> lock(mutex);
        double jitter = distribution(generator);
        return std::chrono::milliseconds(
            static_cast<int>(current_delay.count() * jitter)
        );
    }
    
    void reset() {
        std::lock_guard<std::mutex> lock(mutex);
        current_delay = base_delay;
        request_count = 0;
        std::cout << "🔄 Rate limiter reset to base delay: " << base_delay.count() << "ms" << std::endl;
    }
    
    int get_request_count() const {
        return request_count.load();
    }
    
    std::chrono::milliseconds get_current_delay() const {
        return current_delay;
    }
    
    // OPSEC: Emergency slowdown for suspicious activity
    void emergency_slowdown() {
        std::lock_guard<std::mutex> lock(mutex);
        current_delay = std::chrono::milliseconds(10000); // 10 seconds
        std::cout << "🚨 EMERGENCY SLOWDOWN: Rate limit increased to 10s" << std::endl;
    }
};

#endif // RATE_LIMITER_HPP
