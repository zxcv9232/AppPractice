import requests
from typing import List, Dict, Optional

class APIClient:
    def __init__(self, base_url: str = "http://localhost:8080/api"):
        self.base_url = base_url
    
    def get_prices(self) -> List[Dict]:
        try:
            response = requests.get(f"{self.base_url}/prices", timeout=5)
            response.raise_for_status()
            return response.json().get("data", [])
        except Exception as e:
            print(f"Error fetching prices: {e}")
            return []
    
    def create_alert(self, user_id: str, symbol: str, alert_type: str = "price", 
                     target_price: float = None, direction: str = None,
                     target_volume: float = None, time_window: int = None) -> Optional[Dict]:
        try:
            payload = {
                "userId": user_id,
                "symbol": symbol,
                "alertType": alert_type
            }
            
            if alert_type == "price":
                payload["targetPrice"] = target_price
                payload["direction"] = direction
            elif alert_type == "volume":
                payload["targetVolume"] = target_volume
                payload["timeWindow"] = time_window if time_window else 1
            
            response = requests.post(f"{self.base_url}/alerts", json=payload, timeout=5)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            print(f"Error creating alert: {e}")
            return None
    
    def get_user_alerts(self, user_id: str) -> List[Dict]:
        try:
            response = requests.get(f"{self.base_url}/alerts/{user_id}", timeout=5)
            response.raise_for_status()
            return response.json().get("data", [])
        except Exception as e:
            print(f"Error fetching alerts: {e}")
            return []
    
    def delete_alert(self, alert_id: str) -> bool:
        try:
            response = requests.delete(f"{self.base_url}/alerts/{alert_id}", timeout=5)
            response.raise_for_status()
            return True
        except Exception as e:
            print(f"Error deleting alert: {e}")
            return False

