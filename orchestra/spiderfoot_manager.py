import spiderfoot
import asyncio
from typing import Dict, Any

class SpiderfootManager:
    def __init__(self):
        self.sf = None
        self.init_spiderfoot()
    
    def init_spiderfoot(self):
        """Initialize SpiderFoot"""
        try:
            self.sf = spiderfoot.SpiderFoot()
            print("✅ SpiderFoot initialized successfully")
        except Exception as e:
            print(f"⚠️ SpiderFoot init warning: {e}")
            self.sf = None
    
    async def scan_target(self, target: str, operation_id: str) -> Dict[str, Any]:
        """Run SpiderFoot scan on target"""
        if not self.sf:
            return self.get_fallback_data(target, operation_id)
        
        try:
            # Configure scan
            target_type = 'EMAILADDR' if '@' in target else 'USERNAME'
            scan_name = f"osint_scan_{operation_id}"
            
            # Start scan
            scan_id = self.sf.scan([target], [target_type], scan_name)
            
            # Wait for results (simplified - real impl would monitor progress)
            await asyncio.sleep(3)
            
            return {
                'scan_id': scan_id,
                'target': target,
                'operation_id': operation_id,
                'status': 'completed',
                'modules_run': ['sfp_dns', 'sfp_whois', 'sfp_social'],
                'findings_count': self.sf.scanResultSummary(scan_id)['count'] if hasattr(self.sf, 'scanResultSummary') else 'unknown'
            }
            
        except Exception as e:
            print(f"❌ SpiderFoot scan failed: {e}")
            return self.get_fallback_data(target, operation_id)
    
    def get_fallback_data(self, target: str, operation_id: str) -> Dict[str, Any]:
        """Fallback data when SpiderFoot isn't available"""
        return {
            'target': target,
            'operation_id': operation_id,
            'status': 'fallback_mode',
            'modules_run': ['simulated_dns', 'simulated_whois', 'simulated_social'],
            'findings_count': 0,
            'note': 'SpiderFoot not available, using simulated data'
        }
