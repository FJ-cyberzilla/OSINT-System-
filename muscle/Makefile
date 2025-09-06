# Compiler and flags
CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -pedantic -O2
LDFLAGS = -lcurl -lnlohmann_json::nlohmann_json

# Directories
SRC_DIR = .
BUILD_DIR = build
BIN_DIR = bin

# Targets
TARGET = $(BIN_DIR)/scanner_engine
TEST_TARGET = $(BIN_DIR)/scanner_test

# Source files
SRCS = $(SRC_DIR)/scanner_engine.cpp
TEST_SRCS = $(SRC_DIR)/scanner_engine.test.cpp

# Object files
OBJS = $(BUILD_DIR)/scanner_engine.o
TEST_OBJS = $(BUILD_DIR)/scanner_engine.test.o

# Default target
all: $(TARGET) $(TEST_TARGET)

# Main executable
$(TARGET): $(OBJS) | $(BIN_DIR)
	$(CXX) $(OBJS) -o $@ $(LDFLAGS)

# Test executable
$(TEST_TARGET): $(TEST_OBJS) $(OBJS) | $(BIN_DIR)
	$(CXX) $(TEST_OBJS) $(OBJS) -o $@ $(LDFLAGS)

# Build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Compile source files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp | $(BUILD_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Clean
clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)

# Run tests
test: $(TEST_TARGET)
	./$(TEST_TARGET)

# Run with valgrind for memory checking
valgrind-test: $(TEST_TARGET)
	valgrind --leak-check=full --track-origins=yes ./$(TEST_TARGET)

# Install dependencies (Ubuntu/Debian)
install-deps:
	sudo apt-get update
	sudo apt-get install -y libcurl4-openssl-dev nlohmann-json3-dev make g++

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Build everything (default)"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Build and run tests"
	@echo "  valgrind-test - Run tests with valgrind memory checking"
	@echo "  install-deps - Install required dependencies"

.PHONY: all clean test valgrind-test install-deps help
