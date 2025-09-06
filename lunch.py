#!/usr/bin/env python3
"""
OSINT System Launcher - Termux & Anaconda Compatible
One-command deployment for the complete system
"""

import os
import sys
import platform
import subprocess
import argparse
from pathlib import Path

class OSINTLauncher:
    def __init__(self):
        self.system_type = self.detect_system()
        self.base_dir = Path(__file__).parent
        self.config = self.load_config()
        
    def detect_system(self):
        """Detect if we're running on Termux or Anaconda"""
        if 'com.termux' in os.environ.get('PREFIX', ''):
            return 'termux'
        elif 'anaconda' in sys.executable.lower() or 'conda' in sys.executable.lower():
            return 'anaconda'
        else:
            return 'standard'
    
    def load_config(self):
        """Load system-specific configuration"""
        config_file = self.base_dir / 'config' / f'{self.system_type}_config.json'
        if config_file.exists():
            import json
            return json.loads(config_file.read_text())
        return {}
    
    def install_dependencies(self):
        """Install required dependencies for the detected system"""
        print(f"📦 Installing dependencies for {self.system_type.upper()}...")
        
        if self.system_type == 'termux':
            # Install Termux specific packages
            subprocess.run(['pkg', 'install', '-y', 'python', 'clang', 'make', 'libxml2', 'libxslt'], check=True)
        
        # Install Python dependencies
        subprocess.run([sys.executable, '-m', 'pip', 'install', '-r', 'requirements.txt'], check=True)
        
        # Compile C++ components if needed
        if self.system_type == 'termux':
            self.compile_termux_components()
    
    def compile_termux_components(self):
        """Compile C++ components for Termux"""
        print("🔨 Compiling C++ components for Termux...")
        muscle_dir = self.base_dir / 'muscle'
        if (muscle_dir / 'Makefile').exists():
            subprocess.run(['make', '-C', str(muscle_dir)], check=True)
    
    def setup_directories(self):
        """Create necessary operational directories"""
        directories = ['ops/results', 'ops/logs', 'ops/manifests', 'brain/trained_models']
        for dir_path in directories:
            (self.base_dir / dir_path).mkdir(parents=True, exist_ok=True)
    
    def start_system(self, mode='full'):
        """Start the OSINT system"""
        print(f"🚀 Starting OSINT System in {mode.upper()} mode...")
        
        if mode == 'full':
            # Start all components
            self.start_brain()
            self.start_orchestra()
            self.start_muscle()
        elif mode == 'lightweight':
            # Start only essential components
            self.start_orchestra()
        elif mode == 'ai-only':
            # Start only AI components
            self.start_brain()
    
    def start_brain(self):
        """Start the AI Brain component"""
        print("🧠 Starting AI Brain...")
        brain_script = self.base_dir / 'brain' / 'main.py'
        subprocess.Popen([sys.executable, str(brain_script)])
    
    def start_orchestra(self):
        """Start the Orchestra coordinator"""
        print("🎻 Starting Orchestra...")
        orchestra_script = self.base_dir / 'orchestra' / 'main.py'
        subprocess.Popen([sys.executable, str(orchestra_script)])
    
    def start_muscle(self):
        """Start the Muscle performance layer"""
        print("💪 Starting Muscle...")
        muscle_binary = self.base_dir / 'muscle' / 'scanner_engine'
        if muscle_binary.exists():
            subprocess.Popen([str(muscle_binary)])
        else:
            print("⚠️  Muscle component not compiled. Running in Python mode.")
    
    def run(self):
        """Main launcher routine"""
        parser = argparse.ArgumentParser(description='OSINT System Launcher')
        parser.add_argument('--install', action='store_true', help='Install dependencies')
        parser.add_argument('--mode', choices=['full', 'lightweight', 'ai-only'], 
                          default='full', help='Operation mode')
        parser.add_argument('--target', help='Immediate target to investigate')
        
        args = parser.parse_args()
        
        if args.install:
            self.install_dependencies()
            self.setup_directories()
        
        self.start_system(args.mode)
        
        if args.target:
            # Immediate investigation
            self.run_investigation(args.target)

def main():
    launcher = OSINTLauncher()
    launcher.run()

if __name__ == "__main__":
    main()
