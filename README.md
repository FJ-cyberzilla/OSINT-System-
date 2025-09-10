# 🔎 OSINT System - Professional Intelligence Gathering Platform

⚡ **Warning:** This tool/app/code has been provided for educational purposes only. Use on systems you own or have explicit permission to test.

🎲 A multi-language OSINT (Open-Source Intelligence) system designed for **educational use** in security research and ethical hacking studies.

## ⚠️ LEGAL DISCLAIMER

**USING THIS SOFTWARE FOR EDUCATIONAL PURPOSES ONLY IS MANDATORY USE RESPONSIBLY.**

```
BEFORE USING THIS TOOL:

1. 🚫 DO NOT use on systems without explicit written permission
2. 🚫 DO NOT use for illegal or malicious activities
3. 🚫 DO NOT violate privacy laws or terms of service
4. ✅ DO use only for educational and research purposes
5. ✅ DO ensure compliance with local laws and regulations
6. ✅ DO obtain proper authorization before testing

The authors assume NO liability for misuse of this software.
Use at your own risk and responsibility.
```

A world-class, multi-language OSINT (Open-Source Intelligence) system designed for professional reconnaissance operations. This platform combines AI-powered analysis with military-grade operational security.

## 🚀 Features

### 🤖 AI-Powered Analysis

*   **Pattern Recognition**: ML-based username analysis and classification using models in `brain/AI` (powered by `SecurePatternAnalyzer`).
*   **Behavior Prediction**: Predictive analytics for target behavior.
*   **Risk Assessment**: Intelligent threat scoring with confidence levels (powered by `SecureThreatAssessor`).
*   **Anomaly Detection**: Automated detection of suspicious patterns.

### 🔍 Professional OSINT Tools

*   **Scrapy Integration**: Advanced web crawling capabilities (details in `muscle`).
*   **SpiderFoot Integration**: Comprehensive intelligence gathering.
*   **Multi-Platform Checking**: 50+ social media platforms.
*   **Real-time Monitoring**: Live intelligence feeds.

### 🛡️ Operational Security

*   **Zero Forensic Footprint**: No plaintext data in logs (see `brain/main.py`).
*   **Hashed Operations**: Secure target referencing (SHA256 hashing in `brain/main.py`).
*   **Plausible Deniability**: Clean operational manifests.
*   **Adaptive Rate Limiting**: Anti-detection mechanisms.

### 🏗️ Enterprise Architecture

## 📦 Installation

### Quick Start (Docker)

```bash
# Clone the repository
git clone https://github.com/FJ-cyberzilla/OSINT-System-.git
cd OSINT-System-

# Start all services
docker-compose up -d

# Access the system
curl http://localhost:8080/api/v1/health
```

### Manual Installation

```bash
# 1. Install dependencies
make build-python
make build-go
make build-java
make build-cpp

# 2. Set ENV variables
cp .env.example .env
nano .env

# 3. Start the BRAIN (AI) service
python brain/main.py

# 4. Start the API service
make docker-build
make docker-test
```

### Termux (Android)

```bash
# Run on mobile devices
pkg update && pkg upgrade
pkg install python clang make openjdk-17 go
make termux
```

## 🎯 Usage

### Basic Investigation

```bash
# Single target investigation
curl -X POST http://localhost:8080/api/v1/investigate \
  -H "Content-Type: application/json" \
  -d '{"target": "username", "priority": "high"}'

# Check operation status
curl http://localhost:8080/api/v1/operations/status?id=op_123456
```

### Batch Processing

```bash
# Process multiple targets
echo -e "target1\ntarget2\ntarget3" > targets.txt
python launch.py --batch targets.txt --output json
```

### AI Brain Examples

#### Analyzing a Pattern

This example demonstrates how to use the `analyze_pattern` function to analyze a given target (e.g., a username).

First, you need to send a JSON message to the AI Brain service (running on `tcp://*:5555`):

```json
{
    "action": "analyze_pattern",
    "target": "johndoe123"
}
```

You can send this message using `zmq` or any other suitable tool. Here's an example using `netcat`:

```bash
echo '{"action": "analyze_pattern", "target": "johndoe123"}' | nc localhost 5555
```

The AI Brain service will respond with a JSON message containing the analysis results:

```json
{
    "pattern_type": "human",
    "anomaly_score": -0.25,
    "complexity": 0.7,
    "behavior_profile": "active on social media, interested in technology",
    "confidence": 0.85
}
```

#### Performing a Final Analysis

This example demonstrates how to use the `final_analysis` function to correlate intelligence from multiple sources and assess the overall threat level.

First, you need to gather intelligence data from various sources. This data should be in a structured format (e.g., a list of dictionaries).

Then, you need to send a JSON message to the AI Brain service:

```json
{
    "action": "final_analysis",
    "target": "johndoe123",
    "data_sources": [
        {"source": "Twitter", "activity": "Frequent tweets about cybersecurity"},
        {"source": "LinkedIn", "profile": "Security Analyst at Acme Corp"}
    ]
}
```

The AI Brain service will respond with a JSON message containing the correlated intelligence, threat assessment, predictions, and confidence score:

```json
{
    "correlated_intel": "User is a Security Analyst at Acme Corp and frequently tweets about cybersecurity, indicating a professional interest in the field.",
    "threat_assessment": "Low",
    "predictions": "Likely to attend cybersecurity conferences and engage in related online communities.",
    "confidence_score": 0.9
}
```

### API Integration

```python
import requests

response = requests.post(
    "http://localhost:8080/api/v1/investigate",
    json={"target": "username", "priority": "high"},
    headers={"X-API-Key": "your_api_key"}
)
```

## 🔧 Configuration

### Environment Variables

```bash
# Copy example configuration
cp .env.example .env

# Edit your settings
nano .env
```

### Key Configuration Options

```env
# Security
API_KEY=your_secure_api_key
JWT_SECRET=your_jwt_secret
ENCRYPTION_KEY=your_encryption_key

# Services
BRAIN_PORT=5555
BRAIN_HOST=0.0.0.0
BRAIN_URL=localhost:5555
ORCHESTRATOR_URL=localhost:5558
API_URL=http://localhost:8080

# Operations
LOG_LEVEL=INFO
ENVIRONMENT=production
```

## 🧠 AI Brain Service

The `brain` directory contains the AI/ML logic of the OSINT System. It uses ZeroMQ for inter-process communication and implements OPSEC-compliant logging.

### Components

*   `brain/main.py`: Main entry point for the AI Brain service. Sets up the ZeroMQ socket, configures secure logging, and processes requests.
*   `brain/AI/Machine-Learning-Layer.py`: Implements the core AI/ML logic, including pattern recognition, anomaly detection, and behavioral prediction.
*   `brain/pattern_analysis.py`: Implements secure pattern analysis techniques.
*   `brain/threat_assessment.py`: Implements secure threat assessment techniques.

### Communication

The Brain service communicates with other components using ZeroMQ. It listens for requests on `tcp://*:5555` and responds with JSON messages.

### Logging

The Brain service uses OPSEC-compliant logging. It hashes target information and stores logs in `ops/logs/brain_operations.log`.

## 📊 API Documentation

### Endpoints

*   POST /api/v1/investigate - Start new investigation
*   GET /api/v1/operations/{id} - Get operation status
*   GET /api/v1/health - System health check
*   GET /api/v1/stats - Operational statistics

### Example Response

```json
{
  "operation_id": "op_1700000000_abc123",
  "target": "username",
  "status": "completed",
  "risk_score": 0.72,
  "findings_count": 23,
  "platforms_checked": 15,
  "recommendations": ["Monitor activity", "Verify identity"]
}
```

## 🛡️ Security Features

### OPSEC Compliance

*   No plaintext usernames in logs
*   Hashed target references
*   Secure operational manifests
*   Forensic awareness

### 🧯 Data Protection

*   Encrypted communications
*   Secure storage
*   Clean data handling
*   Regular security audits

## 🛅 Support

*   📖 Documentation
*   💽 Issue Tracker
*   💬 Discussions
*   📧 Email: king.aka.tyrant@hotmail.com

## 🙏 Acknowledgments

*   Scrapy - Web scraping framework
*   SpiderFoot - OSINT automation
*   Spring Boot - Java web framework
*   Docker - Containerization platform
*   GitHub Actions - CI/CD automation
*   ZeroMQ - Inter-process communication

---

[![Security Policy](https://img.shields.io/badge/Security-Policy-blue.svg)](SECURITY.md)

Thanks for your consideration.
FJ™ Cybertronic Systems® - MMXXVI
