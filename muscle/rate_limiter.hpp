// Adaptive Rate Limiting for OPSEC
#ifndef RATE_LIMITER_HPP
#define RATE_LIMITER_HPP

#include <chrono>
#include <random>
#include <atomic>
#include <mutex>

class AdaptiveRateLimiter {
private:
    std::chrono::milliseconds base_delay;
    std::chrono::milliseconds current_delay;
    std::atomic<int> request_count;
    std::mutex mutex;
    std::default_random_engine generator;
    std::uniform_real_distribution<double> distribution;
    
public:
    AdaptiveRateLimiter(int base_delay_ms = 1000) 
        : base_delay(base_delay_ms),
          current_delay(base_delay_ms),
          request_count(0),
          distribution(0.8, 1.2) {}
    
    void wait() {
        std::lock_guard<std::mutex> lock(mutex);
        
        // Add jitter to avoid pattern detection
        double jitter = distribution(generator);
        auto actual_delay = std::chrono::milliseconds(
            static_cast<int>(current_delay.count() * jitter)
        );
        
        std::this_thread::sleep_for(actual_delay);
        
        // Adaptive rate limiting
        request_count++;
        if (request_count % 10 == 0) {
            current_delay = std::chrono::milliseconds(
                std::min(5000, static_cast<int>(current_delay.count() * 1.1))
            );
        }
    }
    
    void reset() {
        std::lock_guard<std::mutex> lock(mutex);
        current_delay = base_delay;
        request_count = 0;
    }
};

#endif // RATE_LIMITER_HPP
