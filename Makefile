# Makefile for Termux C++ compilation
CC = clang++
CFLAGS = -O3 -std=c++17 -Wall
TARGET = scanner_engine
SRCS = main.cpp scanner_engine.cpp rate_limiter.cpp
OBJS = $(SRCS:.cpp=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $(TARGET) $(OBJS)

%.o: %.cpp
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(TARGET) *.o

install: $(TARGET)
	cp $(TARGET) /data/data/com.termux/files/usr/bin/

.PHONY: all clean install
