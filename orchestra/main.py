#!/usr/bin/env python3
"""
ORCHESTRA - Coordination Layer
Coordinates Brain, Muscle, and external components
"""

import zmq
import json
import asyncio
from typing import Dict, Any
from orchestrator import InvestigationOrchestrator

class OrchestraCoordinator:
    def __init__(self):
        self.context = zmq.Context()
        
        # Connect to Brain
        self.brain_socket = self.context.socket(zmq.REQ)
        self.brain_socket.connect("tcp://localhost:5555")
        
        # Connect to Muscle
        self.muscle_socket = self.context.socket(zmq.REQ)
        self.muscle_socket.connect("tcp://localhost:5556")
        
        # Setup server for external connections
        self.server_socket = self.context.socket(zmq.REP)
        self.server_socket.bind("tcp://*:5558")
        
        self.orchestrator = InvestigationOrchestrator()
        
        print("🎻 ORCHESTRA layer initialized and listening on port 5558")
    
    async def coordinate_investigation(self, target: str) -> Dict[str, Any]:
        """Coordinate a complete investigation"""
        return await self.orchestrator.orchestrate(target)
    
    def run(self):
        """Main coordination loop"""
        while True:
            try:
                # Receive request
                message = self.server_socket.recv_json()
                print(f"🎻 Received request: {message.get('action')}")
                
                # Process request based on action
                if message.get('action') == 'investigate':
                    target = message.get('target')
                    
                    # Run investigation asynchronously
                    loop = asyncio.new_event_loop()
                    asyncio.set_event_loop(loop)
                    result = loop.run_until_complete(self.coordinate_investigation(target))
                    loop.close()
                    
                    response = result
                
                else:
                    response = {'error': 'Unknown action'}
                
                # Send response
                self.server_socket.send_json(response)
                
            except Exception as e:
                error_response = {'error': f'Orchestration error: {str(e)}'}
                self.server_socket.send_json(error_response)

if __name__ == "__main__":
    coordinator = OrchestraCoordinator()
    coordinator.run()
