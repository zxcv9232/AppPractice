import flet as ft
import threading
import time
from services.api_client import APIClient

class AlertPage:
    def __init__(self, api_client: APIClient, user_id: str):
        self.api_client = api_client
        self.user_id = user_id
        self.alert_list = ft.ListView(spacing=10, padding=20)
        self.previous_alerts = {}
        self.is_monitoring = False
        self.monitor_thread = None
        self.page = None
        self.home_page = None
        
        self.alert_type_tabs = ft.Tabs(
            selected_index=0,
            tabs=[
                ft.Tab(text="價格警報", icon=ft.Icons.ATTACH_MONEY),
                ft.Tab(text="成交量警報", icon=ft.Icons.SHOW_CHART),
            ],
            on_change=self.on_alert_type_change,
        )
        
        self.symbol_input = ft.TextField(label="幣種代號 (如 BTC)", width=200)
        
        self.price_input = ft.TextField(label="目標價格", width=200, keyboard_type=ft.KeyboardType.NUMBER)
        self.direction_dropdown = ft.Dropdown(
            label="觸發條件",
            width=200,
            options=[
                ft.dropdown.Option("above", "突破此價格"),
                ft.dropdown.Option("below", "跌破此價格"),
            ],
        )
        
        self.volume_input = ft.TextField(label="目標成交量", width=200, keyboard_type=ft.KeyboardType.NUMBER, visible=False)
        self.time_window_input = ft.TextField(label="時間窗口 (分鐘)", width=200, keyboard_type=ft.KeyboardType.NUMBER, value="1", visible=False)
        self.monitoring_status = ft.Container(
            content=ft.Row([
                ft.Icon(ft.Icons.NOTIFICATIONS_ACTIVE, color=ft.Colors.GREEN, size=16),
                ft.Text("全局監控運行中 - 任何頁面都會收到通知", size=12, color=ft.Colors.GREEN)
            ], spacing=5),
            padding=10,
            bgcolor=ft.Colors.GREEN_50,
            border_radius=5,
        )
    
    def set_home_page(self, home_page):
        self.home_page = home_page
    
    def on_alert_type_change(self, e):
        selected_tab = self.alert_type_tabs.selected_index
        
        if selected_tab == 0:
            self.price_input.visible = True
            self.direction_dropdown.visible = True
            self.volume_input.visible = False
            self.time_window_input.visible = False
        else:
            self.price_input.visible = False
            self.direction_dropdown.visible = False
            self.volume_input.visible = True
            self.time_window_input.visible = True
        
        if self.page:
            self.page.update()
    
    def build(self) -> ft.Container:
        return ft.Container(
            content=ft.Column([
                ft.Container(
                    content=ft.Text("警報設定", size=24, weight=ft.FontWeight.BOLD),
                    padding=20,
                ),
                ft.Container(
                    content=self.monitoring_status if self.is_monitoring else None,
                    padding=ft.padding.only(left=20, right=20, bottom=10),
                ),
                ft.Container(
                    content=self.alert_type_tabs,
                    padding=ft.padding.only(left=20, right=20),
                ),
                ft.Container(
                    content=ft.Column([
                        self.symbol_input,
                        self.price_input,
                        self.direction_dropdown,
                        self.volume_input,
                        self.time_window_input,
                        ft.Row([
                            ft.ElevatedButton(
                                text="創建警報",
                                on_click=self.create_alert,
                            ),
                            ft.OutlinedButton(
                                text="🔔 測試通知",
                                on_click=self.test_notification,
                            ),
                        ], spacing=10),
                    ], spacing=10),
                    padding=20,
                ),
                ft.Divider(),
                ft.Container(
                    content=self.alert_list,
                    expand=True,
                ),
            ]),
            expand=True,
        )
    
    def create_alert(self, e):
        symbol = self.symbol_input.value
        alert_type = "price" if self.alert_type_tabs.selected_index == 0 else "volume"
        
        if not symbol:
            self._show_snackbar("請輸入幣種代號", ft.Colors.ORANGE)
            return
        
        if alert_type == "price":
            price = self.price_input.value
            direction = self.direction_dropdown.value
            
            if not all([price, direction]):
                self._show_snackbar("請填寫目標價格和觸發條件", ft.Colors.ORANGE)
                return
            
            result = self.api_client.create_alert(
                user_id=self.user_id,
                symbol=symbol.upper(),
                alert_type="price",
                target_price=float(price),
                direction=direction
            )
            
            if result:
                self.symbol_input.value = ""
                self.price_input.value = ""
                self.direction_dropdown.value = None
        
        elif alert_type == "volume":
            volume = self.volume_input.value
            time_window = self.time_window_input.value
            
            if not volume:
                self._show_snackbar("請填寫目標成交量", ft.Colors.ORANGE)
                return
            
            result = self.api_client.create_alert(
                user_id=self.user_id,
                symbol=symbol.upper(),
                alert_type="volume",
                target_volume=float(volume),
                time_window=int(time_window) if time_window else 1
            )
            
            if result:
                self.symbol_input.value = ""
                self.volume_input.value = ""
                self.time_window_input.value = "1"
        
        if result:
            self.load_alerts(check_triggered=False)
            self._show_snackbar("✅ 警報創建成功！", ft.Colors.GREEN)
    
    def _show_snackbar(self, message: str, bgcolor):
        if self.page:
            snack_bar = ft.SnackBar(
                content=ft.Text(message, size=14),
                bgcolor=bgcolor,
            )
            self.page.overlay.append(snack_bar)
            snack_bar.open = True
            self.page.update()
    
    def load_alerts(self, check_triggered=True):
        alerts = self.api_client.get_user_alerts(self.user_id)
        current_alert_ids = {alert["alertId"] for alert in alerts}
        
        if check_triggered and self.previous_alerts:
            triggered_alerts = []
            for alert_id, alert_info in self.previous_alerts.items():
                if alert_id not in current_alert_ids:
                    triggered_alerts.append(alert_info)
            
            if triggered_alerts:
                print(f"🔔 {len(triggered_alerts)} 個警報已觸發")
                for i, alert_info in enumerate(triggered_alerts):
                    print(f"  {i+1}. {alert_info}")
                self.show_multiple_notifications(triggered_alerts)
        
        self.previous_alerts = {alert["alertId"]: alert for alert in alerts}
        print(f"📊 當前警報數量: {len(alerts)}")
        
        if self.home_page:
            self.home_page.set_alert_count(len(alerts))
        
        self.alert_list.controls.clear()
        
        if not alerts:
            self.alert_list.controls.append(
                ft.Container(
                    content=ft.Text("尚無警報", size=14, color="grey"),
                    padding=20,
                )
            )
        else:
            for alert in alerts:
                alert_type = alert.get("alertType", "price")
                
                if alert_type == "volume":
                    time_window = alert.get("timeWindow", 1)
                    description = ft.Text(
                        f"成交量 {alert['targetVolume']:,.0f} ({time_window}分鐘)", 
                        size=14
                    )
                    icon = ft.Icons.SHOW_CHART
                    icon_color = ft.Colors.BLUE
                else:
                    direction_text = "突破" if alert.get("direction") == "above" else "跌破"
                    description = ft.Text(
                        f"{direction_text} ${alert['targetPrice']:,.2f}", 
                        size=14
                    )
                    icon = ft.Icons.ATTACH_MONEY
                    icon_color = ft.Colors.GREEN if alert.get("direction") == "above" else ft.Colors.RED
                
                card = ft.Card(
                    content=ft.Container(
                        content=ft.Row([
                            ft.Icon(icon, color=icon_color, size=24),
                            ft.Column([
                                ft.Text(f"{alert['symbol']}", size=18, weight=ft.FontWeight.BOLD),
                                description,
                            ], spacing=5),
                            ft.Container(expand=True),
                            ft.IconButton(
                                icon=ft.Icons.DELETE,
                                on_click=lambda e, aid=alert["alertId"]: self.delete_alert(aid),
                            ),
                        ]),
                        padding=15,
                    ),
                )
                self.alert_list.controls.append(card)
        
        if self.alert_list.page:
            self.alert_list.update()
    
    def show_triggered_notification(self, alert_info):
        if not self.page:
            print("⚠️ 警告: page 對象不存在，無法顯示通知")
            return
        
        alert_type = alert_info.get("alertType", "price")
        
        if alert_type == "volume":
            time_window = alert_info.get("timeWindow", 1)
            message = f"🔔 {alert_info['symbol']} 成交量達標！ {alert_info['targetVolume']:,.0f} ({time_window}分鐘)"
            bgcolor = ft.Colors.BLUE
        else:
            direction_text = "突破" if alert_info.get("direction") == "above" else "跌破"
            message = f"🔔 {alert_info['symbol']} 已{direction_text} ${alert_info['targetPrice']:,.2f}!"
            bgcolor = ft.Colors.GREEN if alert_info.get("direction") == "above" else ft.Colors.RED
        
        print(f"✨ 顯示通知: {message}")
        
        snack_bar = ft.SnackBar(
            content=ft.Text(message, size=16, weight=ft.FontWeight.BOLD),
            bgcolor=bgcolor,
            duration=5000,
        )
        
        self.page.overlay.append(snack_bar)
        snack_bar.open = True
        self.page.update()
    
    def show_multiple_notifications(self, triggered_alerts):
        if not self.page or not triggered_alerts:
            return
        
        if len(triggered_alerts) == 1:
            self.show_triggered_notification(triggered_alerts[0])
        else:
            messages = []
            for alert in triggered_alerts:
                alert_type = alert.get("alertType", "price")
                
                if alert_type == "volume":
                    time_window = alert.get("timeWindow", 1)
                    messages.append(f"{alert['symbol']} 成交量 {alert['targetVolume']:,.0f} ({time_window}分)")
                else:
                    direction_text = "突破" if alert.get("direction") == "above" else "跌破"
                    messages.append(f"{alert['symbol']} {direction_text} ${alert['targetPrice']:,.2f}")
            
            combined_message = f"🔔 {len(triggered_alerts)} 個警報已觸發:\n" + "\n".join(f"• {msg}" for msg in messages)
            
            print(f"✨ 顯示多個通知: {combined_message}")
            
            snack_bar = ft.SnackBar(
                content=ft.Text(combined_message, size=14, weight=ft.FontWeight.BOLD),
                bgcolor=ft.Colors.BLUE,
                duration=8000,
            )
            
            self.page.overlay.append(snack_bar)
            snack_bar.open = True
            self.page.update()
    
    def test_notification(self, e):
        if not self.page:
            print("⚠️ 警告: page 對象不存在")
            return
        
        test_alert = {
            "symbol": "BTC",
            "direction": "above",
            "targetPrice": 87000.00
        }
        
        print("🧪 測試通知功能...")
        self.show_triggered_notification(test_alert)
    
    def start_monitoring(self, page):
        if self.is_monitoring:
            print("⚠️ 警報監控已在運行中")
            return
        
        self.page = page
        self.is_monitoring = True
        print("🎯 全局警報監控已啟動 - 在任何頁面都會收到通知")
        
        def monitor_loop():
            while self.is_monitoring:
                time.sleep(5)
                if self.is_monitoring:
                    try:
                        self.load_alerts(check_triggered=True)
                    except Exception as e:
                        print(f"❌ 警報監控錯誤: {e}")
        
        self.monitor_thread = threading.Thread(target=monitor_loop, daemon=True)
        self.monitor_thread.start()
    
    def stop_monitoring(self):
        self.is_monitoring = False
    
    def delete_alert(self, alert_id: str):
        if self.api_client.delete_alert(alert_id):
            self.load_alerts()

