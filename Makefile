# OSINT System Build Automation
SHELL := /bin/bash

# Configuration
PROJECT_NAME = osint-system
VERSION = 1.0.0

# Components
COMPONENTS = brain orchestra api web muscle

# Docker
DOCKER_COMPOSE = docker-compose
DOCKER_BUILD = docker build

# Python
PYTHON = python3
PIP = pip3

# Go
GO = go

# Java
MAVEN = mvn

.PHONY: all test build clean deploy help

# Default target
all: test build

# === TESTING ===
test: test-python test-go test-java test-cpp test-integration

test-python:
	@echo "Testing Python components..."
	cd brain && $(PYTHON) -m pytest tests/ -v --cov=. --cov-report=xml
	cd orchestra && $(PYTHON) -m pytest tests/ -v --cov=. --cov-report=xml

test-go:
	@echo "Testing Go API..."
	cd api && $(GO) test ./... -v -coverprofile=coverage.out

test-java:
	@echo "Testing Java Web..."
	cd web && $(MAVEN) test -B

test-cpp:
	@echo "Testing C++ Muscle..."
	cd muscle && make test

test-integration:
	@echo "Running integration tests..."
	$(DOCKER_COMPOSE) -f docker-compose.test.yml up -d --build
	sleep 20
	curl -f http://localhost:8080/api/v1/health || (echo "Integration test failed" && exit 1)
	$(DOCKER_COMPOSE) -f docker-compose.test.yml down
	@echo "Integration tests passed!"

# === BUILDING ===
build: build-python build-go build-java build-cpp

build-python:
	@echo "Building Python components..."
	cd brain && $(PIP) install -r requirements.txt
	cd orchestra && $(PIP) install -r requirements.txt

build-go:
	@echo "Building Go API..."
	cd api && $(GO) mod download && $(GO) build -v

build-java:
	@echo "Building Java Web..."
	cd web && $(MAVEN) clean package -DskipTests

build-cpp:
	@echo "Building C++ Muscle..."
	cd muscle && make

# === DOCKER ===
docker-build:
	@echo "Building Docker images..."
	for component in $(COMPONENTS); do \
		if [ -f "$$component/Dockerfile" ]; then \
			echo "Building $$component..."; \
			cd $$component && $(DOCKER_BUILD) -t $(PROJECT_NAME)-$$component:$(VERSION) . ; \
			cd ..; \
		fi; \
	done

docker-test:
	@echo "Testing Docker images..."
	for component in $(COMPONENTS); do \
		if [ -f "$$component/Dockerfile" ]; then \
			echo "Testing $$component image..."; \
			docker run --rm $(PROJECT_NAME)-$$component:$(VERSION) echo "$$component image test passed"; \
		fi; \
	done

# === TERMUX ===
termux:
	@echo "Building for Termux..."
	# Simulate Termux environment setup
	cd muscle && make CC="clang++" CFLAGS="-O3 -std=c++17 -Wall"
	cd api && $(GO) build -v
	cd brain && $(PIP) install -r requirements.txt --user
	cd orchestra && $(PIP) install -r requirements.txt --user

# === DEPLOYMENT ===
deploy-staging:
	@echo "Deploying to staging..."
	$(DOCKER_COMPOSE) -f docker-compose.staging.yml up -d --build

deploy-production:
	@echo "Deploying to production..."
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml up -d --build

# === SECURITY ===
security-audit:
	@echo "Running security audit..."
	$(PIP) install safety
	safety check -r brain/requirements.txt
	safety check -r orchestra/requirements.txt
	cd api && $(GO) vet ./...
	cd web && $(MAVEN) dependency-check:check

# === CLEANUP ===
clean:
	@echo "Cleaning build artifacts..."
	for component in $(COMPONENTS); do \
		if [ -d "$$component" ]; then \
			echo "Cleaning $$component..."; \
			cd $$component; \
			if [ -f "Makefile" ]; then make clean; fi; \
			rm -rf __pycache__ *.pyc *.pyo .pytest_cache 2>/dev/null || true; \
			cd ..; \
		fi; \
	done
	rm -rf ops/logs/* ops/results/* ops/temp/*
	$(DOCKER_COMPOSE) down --rmi all --volumes 2>/dev/null || true

# === HELP ===
help:
	@echo "OSINT System Makefile Commands:"
	@echo "  make all              - Run tests and build"
	@echo "  make test             - Run all tests"
	@echo "  make test-python      - Test Python components"
	@echo "  make test-go          - Test Go API"
	@echo "  make test-java        - Test Java Web"
	@echo "  make test-cpp         - Test C++ Muscle"
	@echo "  make build            - Build all components"
	@echo "  make docker-build     - Build Docker images"
	@echo "  make docker-test      - Test Docker images"
	@echo "  make termux           - Build for Termux"
	@echo "  make deploy-staging   - Deploy to staging"
	@echo "  make deploy-production - Deploy to production"
	@echo "  make security-audit   - Run security audit"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make help             - Show this help message"
