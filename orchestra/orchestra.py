import asyncio
from scrapy_manager import ScrapyManager
from spiderfoot_manager import SpiderfootManager

class Orchestrator:
    def __init__(self):
        self.scrapy_manager = ScrapyManager()
        self.spiderfoot_manager = SpiderfootManager()
        self.operation_cache = {}
    
    async def investigate_target(self, target: str, operation_id: str) -> Dict[str, Any]:
        """Complete investigation using all integrated tools"""
        
        # Run all tools in parallel
        scrapy_task = asyncio.to_thread(
            self.scrapy_manager.crawl_target, target, operation_id
        )
        
        spiderfoot_task = self.spiderfoot_manager.scan_target(target, operation_id)
        
        # Wait for both to complete
        scrapy_results, spiderfoot_results = await asyncio.gather(
            scrapy_task, spiderfoot_task
        )
        
        # Correlate results
        final_report = self.correlate_intelligence(
            target, operation_id, scrapy_results, spiderfoot_results
        )
        
        return final_report
    
    def correlate_intelligence(self, target, operation_id, scrapy_data, spiderfoot_data):
        """Correlate data from all sources"""
        return {
            'operation_id': operation_id,
            'target': target,
            'scrapy': {
                'pages_crawled': len(scrapy_data),
                'findings': scrapy_data[:5]  # First 5 results
            },
            'spiderfoot': {
                'scan_id': spiderfoot_data.get('scan_id'),
                'findings_count': spiderfoot_data.get('findings_count', 0),
                'status': spiderfoot_data.get('status')
            },
            'correlation': {
                'risk_score': self.calculate_risk_score(scrapy_data, spiderfoot_data),
                'confidence': 0.85,
                'timestamp': self.get_timestamp()
            }
        }
