import numpy as np
from sklearn.ensemble import RandomForestClassifier
import joblib
import hashlib

class PatternAnalyzer:
    def __init__(self):
        self.model = self.load_model()
        
    def load_model(self):
        """Load or create pattern recognition model"""
        try:
            return joblib.load('brain/trained_models/pattern_classifier.joblib')
        except:
            # Create and train a simple model
            model = RandomForestClassifier(n_estimators=50)
            # This would be trained on real data in production
            return model
    
    def extract_features(self, username):
        """Extract features from username for ML analysis"""
        features = [
            len(username),
            sum(1 for c in username if c.isdigit()),
            sum(1 for c in username if c.isalpha()),
            sum(1 for c in username if not c.isalnum()),
            username.count('_'),
            username.count('.'),
            username.count('-'),
        ]
        return np.array([features])
    
    def analyze(self, username):
        """Analyze username pattern"""
        features = self.extract_features(username)
        prediction = self.model.predict(features)[0]
        
        return {
            'pattern_type': ['generic', 'personal', 'professional', 'automated'][prediction],
            'confidence': 0.85,
            'features': features.tolist(),
            'risk_score': self.calculate_risk_score(features)
        }
    
    def calculate_risk_score(self, features):
        """Calculate risk score based on features"""
        return min(1.0, features[0][1] * 0.1 + features[0][4] * 0.2)  # Simplified
