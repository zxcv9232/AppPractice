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

    # ==================== 指標監控 API ====================
    
    def create_indicator_subscription(
        self, 
        user_id: str, 
        symbol: str, 
        telegram_chat_id: str,
        notify_interval_min: int = 60,
        enable_volume_check: bool = False,
        volume_check_mode: str = "multiplier",
        volume_fixed_value: float = 0,
        volume_multiplier: float = 2.0,
        volume_avg_period: int = 20
    ) -> Optional[Dict]:
        """創建指標監控訂閱"""
        try:
            payload = {
                "userId": user_id,
                "symbol": symbol,
                "telegramChatId": telegram_chat_id,
                "notifyIntervalMin": notify_interval_min,
                "enableVolumeCheck": enable_volume_check,
                "volumeCheckMode": volume_check_mode,
                "volumeFixedValue": volume_fixed_value,
                "volumeMultiplier": volume_multiplier,
                "volumeAvgPeriod": volume_avg_period
            }
            
            response = requests.post(f"{self.base_url}/indicators/subscribe", json=payload, timeout=5)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            print(f"Error creating indicator subscription: {e}")
            return None
    
    def get_indicator_subscriptions(self, user_id: str) -> List[Dict]:
        """獲取用戶的指標訂閱列表"""
        try:
            response = requests.get(f"{self.base_url}/indicators/subscriptions", params={"userId": user_id}, timeout=5)
            response.raise_for_status()
            return response.json() if response.json() else []
        except Exception as e:
            print(f"Error fetching indicator subscriptions: {e}")
            return []
    
    def update_indicator_subscription(self, subscription_id: str, **kwargs) -> Optional[Dict]:
        """更新指標訂閱"""
        try:
            response = requests.put(f"{self.base_url}/indicators/subscriptions/{subscription_id}", json=kwargs, timeout=5)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            print(f"Error updating indicator subscription: {e}")
            return None
    
    def toggle_indicator_subscription(self, subscription_id: str) -> Optional[Dict]:
        """切換訂閱開關"""
        try:
            response = requests.post(f"{self.base_url}/indicators/subscriptions/{subscription_id}/toggle", timeout=5)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            print(f"Error toggling indicator subscription: {e}")
            return None
    
    def delete_indicator_subscription(self, subscription_id: str) -> bool:
        """刪除指標訂閱"""
        try:
            response = requests.delete(f"{self.base_url}/indicators/subscriptions/{subscription_id}", timeout=5)
            response.raise_for_status()
            return True
        except Exception as e:
            print(f"Error deleting indicator subscription: {e}")
            return False
    
    def get_indicator_result(self, symbol: str) -> Optional[Dict]:
        """獲取幣種的指標計算結果"""
        try:
            response = requests.get(f"{self.base_url}/indicators/{symbol}", timeout=5)
            response.raise_for_status()
            return response.json()
        except Exception as e:
            print(f"Error fetching indicator result: {e}")
            return None

