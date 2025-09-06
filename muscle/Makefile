# Muscle Component Makefile
CC = g++
CFLAGS = -O3 -std=c++17 -Wall -Wextra -I.
TARGET = scanner_engine
SRCS = $(wildcard *.cpp)
OBJS = $(SRCS:.cpp=.o)
TEST_TARGET = muscle_test
TEST_SRCS = $(wildcard tests/*.cpp)

.PHONY: all clean test

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJS)

%.o: %.cpp
	$(CC) $(CFLAGS) -c $< -o $@

test: $(TEST_TARGET)
	./$(TEST_TARGET)

$(TEST_TARGET): $(filter-out main.o, $(OBJS)) $(TEST_SRCS:.cpp=.o)
	$(CC) $(CFLAGS) -o $(TEST_TARGET) $^

tests/%.o: tests/%.cpp
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET) $(TEST_TARGET) *.o tests/*.o

install: $(TARGET)
	cp $(TARGET) /usr/local/bin/
