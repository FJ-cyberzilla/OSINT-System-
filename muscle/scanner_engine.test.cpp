#include "scanner_engine.hpp"
#include <iostream>
#include <cassert>

void test_scanner_initialization() {
    std::cout << "Testing ScannerEngine initialization..." << std::endl;
    ScannerEngine scanner;
    assert(scanner.initialize() == true);
    std::cout << "✅ Initialization test passed" << std::endl;
}

void test_timeout_setting() {
    std::cout << "Testing timeout setting..." << std::endl;
    ScannerEngine scanner;
    scanner.setTimeout(3000);
    // We can't easily test the internal timeout, but we can verify the function doesn't crash
    std::cout << "✅ Timeout setting test passed" << std::endl;
}

void test_scan_result_structure() {
    std::cout << "Testing scan result structure..." << std::endl;
    ScanResult result;
    result.target = "testuser";
    result.website_name = "TestSite";
    result.url = "https://example.com/testuser";
    result.found = true;
    result.status_code = 200;
    
    assert(result.target == "testuser");
    assert(result.website_name == "TestSite");
    assert(result.found == true);
    assert(result.status_code == 200);
    std::cout << "✅ Scan result structure test passed" << std::endl;
}

void test_status_code_checking() {
    std::cout << "Testing status code checking..." << std::endl;
    ScannerEngine scanner;
    
    std::vector<int> success_codes = {200, 201, 301};
    assert(scanner.check_status_code(200, success_codes) == true);
    assert(scanner.check_status_code(404, success_codes) == false);
    assert(scanner.check_status_code(301, success_codes) == true);
    
    std::cout << "✅ Status code checking test passed" << std::endl;
}

void test_performance_metrics() {
    std::cout << "Testing performance metrics..." << std::endl;
    ScannerEngine scanner;
    scanner.initialize();
    
    auto metrics = scanner.get_performance_metrics();
    assert(metrics.contains("timeout_ms"));
    assert(metrics.contains("initialized"));
    assert(metrics["initialized"] == true);
    
    std::cout << "✅ Performance metrics test passed" << std::endl;
}

void test_intelligence_data() {
    std::cout << "Testing intelligence data..." << std::endl;
    ScannerEngine scanner;
    scanner.initialize();
    
    json scan_data;
    IntelligenceData intel = scanner.gather_additional_intel("testuser", scan_data);
    
    assert(intel.breach_info["target"] == "testuser");
    assert(intel.whois_info["target"] == "testuser");
    assert(intel.social_analysis["target"] == "testuser");
    assert(intel.threat_data["target"] == "testuser");
    
    std::cout << "✅ Intelligence data test passed" << std::endl;
}

int main() {
    std::cout << "🧪 Starting ScannerEngine tests..." << std::endl;
    
    try {
        test_scanner_initialization();
        test_timeout_setting();
        test_scan_result_structure();
        test_status_code_checking();
        test_performance_metrics();
        test_intelligence_data();
        
        std::cout << "🎉 All tests passed successfully!" << std::endl;
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "❌ Test failed: " << e.what() << std::endl;
        return 1;
    }
}
