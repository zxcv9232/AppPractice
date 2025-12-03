import flet as ft
import threading
import time
from services.api_client import APIClient
from utils.helpers import format_price, format_change, get_change_color

class HomePage:
    def __init__(self, api_client: APIClient, on_symbol_click):
        self.api_client = api_client
        self.on_symbol_click = on_symbol_click
        self.price_list = ft.ListView(spacing=10, padding=20)
        self.refresh_timer = None
        self.is_active = False
        self.last_update_text = ft.Text("", size=12, color="grey")
        self.alert_count_badge = None
    
    def build(self) -> ft.Container:
        self.alert_count_badge = ft.Container(
            content=ft.Row([
                ft.Icon(ft.Icons.NOTIFICATIONS_ACTIVE, color=ft.Colors.ORANGE, size=14),
                ft.Text("警報監控中", size=11, color=ft.Colors.ORANGE)
            ], spacing=3),
            padding=5,
            bgcolor=ft.Colors.ORANGE_50,
            border_radius=3,
            visible=False,
        )
        
        return ft.Container(
            content=ft.Column([
                ft.Container(
                    content=ft.Column([
                        ft.Row([
                            ft.Text("虛擬幣看盤", size=24, weight=ft.FontWeight.BOLD),
                            ft.Container(expand=True),
                            self.last_update_text,
                        ]),
                        self.alert_count_badge,
                    ], spacing=5),
                    padding=20,
                ),
                ft.Container(
                    content=self.price_list,
                    expand=True,
                ),
            ]),
            expand=True,
        )
    
    def set_alert_count(self, count: int):
        if self.alert_count_badge:
            if count > 0:
                self.alert_count_badge.content.controls[1].value = f"{count} 個警報監控中"
                self.alert_count_badge.visible = True
            else:
                self.alert_count_badge.visible = False
            if self.alert_count_badge.page:
                self.alert_count_badge.update()
    
    def load_prices(self):
        prices = self.api_client.get_prices()
        self.price_list.controls.clear()
        
        if not prices:
            self.price_list.controls.append(
                ft.Container(
                    content=ft.Column([
                        ft.Text("⚠️ 無法連接到後端 API", size=16, color="red"),
                        ft.Text("請確認後端服務正在運行", size=12),
                        ft.Text("cd backend && go run cmd/server/main.go", 
                               size=10, color="grey"),
                    ], horizontal_alignment=ft.CrossAxisAlignment.CENTER),
                    padding=20,
                )
            )
        else:
            for price_data in prices:
                symbol = price_data["symbol"]
                price = price_data["price"]
                change = price_data["change24h"]
                
                card = ft.Card(
                    content=ft.Container(
                        content=ft.Row([
                            ft.Column([
                                ft.Text(symbol, size=20, weight=ft.FontWeight.BOLD),
                                ft.Text(format_price(price), size=16),
                            ], spacing=5),
                            ft.Container(expand=True),
                            ft.Text(
                                format_change(change),
                                size=16,
                                color=get_change_color(change),
                                weight=ft.FontWeight.BOLD,
                            ),
                        ]),
                        padding=15,
                        on_click=lambda e, s=symbol: self.on_symbol_click(s),
                    ),
                )
                self.price_list.controls.append(card)
        
        current_time = time.strftime("%H:%M:%S")
        self.last_update_text.value = f"最後更新: {current_time}"
        
        if self.price_list.page:
            self.price_list.update()
            if self.last_update_text.page:
                self.last_update_text.update()
    
    def start_auto_refresh(self):
        self.is_active = True
        
        def refresh_loop():
            while self.is_active:
                time.sleep(10)
                if self.is_active and self.price_list.page:
                    try:
                        self.load_prices()
                    except Exception as e:
                        print(f"自動刷新錯誤: {e}")
        
        self.refresh_timer = threading.Thread(target=refresh_loop, daemon=True)
        self.refresh_timer.start()
    
    def stop_auto_refresh(self):
        self.is_active = False

