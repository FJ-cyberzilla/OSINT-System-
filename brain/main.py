#!/usr/bin/env python3
"""
BRAIN - AI/ML Layer with OPSEC Features
"""

import zmq
import json
import hashlib
import logging
from datetime import datetime
from pattern_analysis import SecurePatternAnalyzer
from threat_assessment import SecureThreatAssessor

class SecureAIBrain:
    def __init__(self):
        self.context = zmq.Context()
        self.socket = self.context.socket(zmq.REP)
        self.socket.bind("tcp://*:5555")
        
        # Secure AI components
        self.pattern_analyzer = SecurePatternAnalyzer()
        self.threat_assessor = SecureThreatAssessor()
        
        # Configure secure logging
        self.setup_secure_logging()
        
        logging.info("🧠 Secure BRAIN layer initialized")

    def setup_secure_logging(self):
        """Configure OPSEC-compliant logging"""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - OPERATION[%(operation_id)s] - %(message)s',
            handlers=[
                logging.FileHandler('ops/logs/brain_operations.log'),
                logging.StreamHandler()
            ]
        )

    def hash_target(self, target: str) -> str:
        """Hash target for OPSEC security"""
        return hashlib.sha256(target.encode()).hexdigest()[:16]

    def process_request(self, message: dict) -> dict:
        """Process request with OPSEC measures"""
        operation_id = message.get('operation_id', 'unknown')
        target_hash = self.hash_target(message.get('target', ''))
        
        logging.info(f"Processing request for target hash: {target_hash}")
        
        action = message.get('action')
        if action == 'analyze_pattern':
            return self.pattern_analyzer.analyze_secure(
                message['target'], 
                operation_id
            )
        elif action == 'assess_threat':
            return self.threat_assessor.assess_secure(
                message['intelligence_data'],
                operation_id
            )
        else:
            return {'error': 'Unknown action', 'operation_id': operation_id}

    def run(self):
        """Secure main processing loop"""
        while True:
            try:
                message = self.socket.recv_json()
                response = self.process_request(message)
                self.socket.send_json(response)
            except Exception as e:
                error_msg = f"Processing error: {str(e)}"
                logging.error(error_msg)
                self.socket.send_json({'error': error_msg})

if __name__ == "__main__":
    brain = SecureAIBrain()
    brain.run()
