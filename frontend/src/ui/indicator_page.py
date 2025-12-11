import flet as ft
from services.api_client import APIClient

class IndicatorPage:
    """指標監控頁面 - LRC 線性回歸通道突破警報"""
    
    def __init__(self, api_client: APIClient, user_id: str):
        self.api_client = api_client
        self.user_id = user_id
        self.page = None
        self.subscription_list = ft.ListView(spacing=10, padding=20)
        
        # 輸入欄位
        self.telegram_chat_id_input = ft.TextField(
            label="Telegram Chat ID",
            hint_text="從 @TradeApocalypse_bot 獲取",
            width=280,
            prefix_icon=ft.Icons.TELEGRAM,
        )
        
        self.symbol_dropdown = ft.Dropdown(
            label="選擇幣種",
            width=280,
            options=[
                ft.dropdown.Option("BTC", "BTC - 比特幣"),
                ft.dropdown.Option("ETH", "ETH - 以太坊"),
                ft.dropdown.Option("BNB", "BNB - 幣安幣"),
                ft.dropdown.Option("SOL", "SOL - Solana"),
                ft.dropdown.Option("XRP", "XRP - 瑞波幣"),
                ft.dropdown.Option("DOGE", "DOGE - 狗狗幣"),
                ft.dropdown.Option("ADA", "ADA - Cardano"),
                ft.dropdown.Option("AVAX", "AVAX - Avalanche"),
                ft.dropdown.Option("1000SHIB", "SHIB - Shiba Inu"),
                ft.dropdown.Option("BCH", "BCH - Bitcoin Cash"),
                ft.dropdown.Option("DOT", "DOT - Polkadot"),
                ft.dropdown.Option("LINK", "LINK - Chainlink"),
                ft.dropdown.Option("TON", "TON - Toncoin"),
                ft.dropdown.Option("UNI", "UNI - Uniswap"),
                ft.dropdown.Option("LTC", "LTC - Litecoin"),
                ft.dropdown.Option("NEAR", "NEAR - NEAR Protocol"),
                ft.dropdown.Option("ATOM", "ATOM - Cosmos"),
                ft.dropdown.Option("AAVE", "AAVE - Aave"),
                ft.dropdown.Option("RIVER", "RIVER - River"),
            ],
        )
        
        self.notify_interval_dropdown = ft.Dropdown(
            label="通知間隔",
            width=280,
            value="60",
            options=[
                ft.dropdown.Option("30", "30 分鐘"),
                ft.dropdown.Option("60", "1 小時"),
                ft.dropdown.Option("120", "2 小時"),
                ft.dropdown.Option("240", "4 小時"),
            ],
        )
        
        # 成交量設定
        self.enable_volume_check = ft.Switch(
            label="啟用成交量判斷",
            value=False,
            on_change=self.on_volume_check_change,
        )
        
        self.volume_mode_dropdown = ft.Dropdown(
            label="成交量模式",
            width=280,
            value="multiplier",
            visible=False,
            options=[
                ft.dropdown.Option("multiplier", "倍數模式 (N 倍均量)"),
                ft.dropdown.Option("fixed", "固定值模式"),
            ],
            on_change=self.on_volume_mode_change,
        )
        
        self.volume_multiplier_input = ft.TextField(
            label="成交量倍數",
            hint_text="例如: 2.0 表示 2 倍均量",
            width=280,
            value="2.0",
            visible=False,
            keyboard_type=ft.KeyboardType.NUMBER,
        )
        
        self.volume_fixed_input = ft.TextField(
            label="固定成交量閾值",
            hint_text="例如: 1000",
            width=280,
            visible=False,
            keyboard_type=ft.KeyboardType.NUMBER,
        )
        
        # 指標結果顯示區
        self.indicator_result_container = ft.Container(
            content=ft.Text("選擇幣種後可查看當前指標", color=ft.Colors.GREY),
            padding=15,
            bgcolor=ft.Colors.GREY_100,
            border_radius=10,
        )
    
    def set_page(self, page):
        self.page = page
    
    def on_volume_check_change(self, e):
        self.volume_mode_dropdown.visible = e.control.value
        self.volume_multiplier_input.visible = e.control.value and self.volume_mode_dropdown.value == "multiplier"
        self.volume_fixed_input.visible = e.control.value and self.volume_mode_dropdown.value == "fixed"
        if self.page:
            self.page.update()
    
    def on_volume_mode_change(self, e):
        self.volume_multiplier_input.visible = e.control.value == "multiplier"
        self.volume_fixed_input.visible = e.control.value == "fixed"
        if self.page:
            self.page.update()
    
    def build(self) -> ft.Container:
        return ft.Container(
            content=ft.Column([
                # 標題
                ft.Container(
                    content=ft.Column([
                        ft.Text("📊 LRC 指標監控", size=24, weight=ft.FontWeight.BOLD),
                        ft.Text("U本位永續合約 - 線性回歸通道突破警報", size=12, color=ft.Colors.GREY),
                    ]),
                    padding=20,
                ),
                
                # 說明卡片
                ft.Container(
                    content=ft.Card(
                        content=ft.Container(
                            content=ft.Column([
                                ft.Row([
                                    ft.Icon(ft.Icons.INFO_OUTLINE, color=ft.Colors.BLUE),
                                    ft.Text("使用說明", weight=ft.FontWeight.BOLD),
                                ]),
                                ft.Text(
                                    "1. 在 Telegram 搜尋 @TradeApocalypse_bot\n"
                                    "2. 發送 /start 獲取你的 Chat ID\n"
                                    "3. 在下方輸入 Chat ID 並選擇要監控的幣種\n"
                                    "4. 監控 U本位永續合約 4H LRC 指標\n"
                                    "5. 當價格突破上/下軌時，Telegram 會收到通知",
                                    size=12,
                                ),
                            ], spacing=5),
                            padding=15,
                        ),
                    ),
                    padding=ft.padding.only(left=20, right=20),
                ),
                
                # 輸入表單
                ft.Container(
                    content=ft.Column([
                        self.telegram_chat_id_input,
                        self.symbol_dropdown,
                        self.notify_interval_dropdown,
                        ft.Divider(),
                        self.enable_volume_check,
                        self.volume_mode_dropdown,
                        self.volume_multiplier_input,
                        self.volume_fixed_input,
                        ft.Row([
                            ft.ElevatedButton(
                                text="創建訂閱",
                                icon=ft.Icons.ADD_ALERT,
                                on_click=self.create_subscription,
                            ),
                            ft.OutlinedButton(
                                text="查看指標",
                                icon=ft.Icons.SHOW_CHART,
                                on_click=self.view_indicator,
                            ),
                        ], spacing=10),
                    ], spacing=10),
                    padding=20,
                ),
                
                # 指標結果
                ft.Container(
                    content=self.indicator_result_container,
                    padding=ft.padding.only(left=20, right=20),
                ),
                
                ft.Divider(),
                
                # 訂閱列表標題
                ft.Container(
                    content=ft.Text("我的訂閱", size=18, weight=ft.FontWeight.BOLD),
                    padding=ft.padding.only(left=20, top=10),
                ),
                
                # 訂閱列表
                ft.Container(
                    content=self.subscription_list,
                    expand=True,
                ),
            ], scroll=ft.ScrollMode.AUTO),
            expand=True,
        )
    
    def create_subscription(self, e):
        telegram_chat_id = self.telegram_chat_id_input.value
        symbol = self.symbol_dropdown.value
        notify_interval = int(self.notify_interval_dropdown.value or "60")
        
        if not telegram_chat_id:
            self._show_snackbar("請輸入 Telegram Chat ID", ft.Colors.ORANGE)
            return
        
        if not symbol:
            self._show_snackbar("請選擇幣種", ft.Colors.ORANGE)
            return
        
        # 成交量設定
        enable_volume = self.enable_volume_check.value
        volume_mode = self.volume_mode_dropdown.value
        volume_multiplier = float(self.volume_multiplier_input.value or "2.0")
        volume_fixed = float(self.volume_fixed_input.value or "0")
        
        result = self.api_client.create_indicator_subscription(
            user_id=self.user_id,
            symbol=symbol,
            telegram_chat_id=telegram_chat_id,
            notify_interval_min=notify_interval,
            enable_volume_check=enable_volume,
            volume_check_mode=volume_mode,
            volume_fixed_value=volume_fixed,
            volume_multiplier=volume_multiplier,
        )
        
        if result:
            self._show_snackbar("✅ 訂閱創建成功！", ft.Colors.GREEN)
            self.symbol_dropdown.value = None
            self.load_subscriptions()
        else:
            self._show_snackbar("❌ 訂閱創建失敗", ft.Colors.RED)
        
        if self.page:
            self.page.update()
    
    def view_indicator(self, e):
        symbol = self.symbol_dropdown.value
        if not symbol:
            self._show_snackbar("請先選擇幣種", ft.Colors.ORANGE)
            return
        
        result = self.api_client.get_indicator_result(symbol)
        if result:
            # 格式化顯示
            price_status = "🔴 跌破下軌" if result.get("isBelowLower") else ("🟢 突破上軌" if result.get("isAboveUpper") else "⚪ 在通道內")
            
            self.indicator_result_container.content = ft.Column([
                ft.Row([
                    ft.Text(f"{symbol}", size=20, weight=ft.FontWeight.BOLD),
                    ft.Container(
                        content=ft.Text(price_status, size=12),
                        bgcolor=ft.Colors.RED_100 if result.get("isBelowLower") else (ft.Colors.GREEN_100 if result.get("isAboveUpper") else ft.Colors.GREY_200),
                        padding=5,
                        border_radius=5,
                    ),
                ], alignment=ft.MainAxisAlignment.SPACE_BETWEEN),
                ft.Divider(),
                ft.Row([
                    ft.Column([
                        ft.Text("當前價格", size=10, color=ft.Colors.GREY),
                        ft.Text(f"${result.get('currentPrice', 0):,.2f}", size=16, weight=ft.FontWeight.BOLD),
                    ]),
                    ft.Column([
                        ft.Text("上軌", size=10, color=ft.Colors.GREY),
                        ft.Text(f"${result.get('upperBand', 0):,.2f}", size=14, color=ft.Colors.GREEN),
                    ]),
                    ft.Column([
                        ft.Text("下軌", size=10, color=ft.Colors.GREY),
                        ft.Text(f"${result.get('lowerBand', 0):,.2f}", size=14, color=ft.Colors.RED),
                    ]),
                ], alignment=ft.MainAxisAlignment.SPACE_AROUND),
                ft.Row([
                    ft.Column([
                        ft.Text("中線", size=10, color=ft.Colors.GREY),
                        ft.Text(f"${result.get('centerLine', 0):,.2f}", size=12),
                    ]),
                    ft.Column([
                        ft.Text("成交量", size=10, color=ft.Colors.GREY),
                        ft.Text(f"{result.get('currentVolume', 0):,.0f}", size=12),
                    ]),
                    ft.Column([
                        ft.Text("量比", size=10, color=ft.Colors.GREY),
                        ft.Text(f"{result.get('volumeRatio', 0):.2f}x", size=12),
                    ]),
                ], alignment=ft.MainAxisAlignment.SPACE_AROUND),
            ], spacing=10)
        else:
            self.indicator_result_container.content = ft.Text("無法獲取指標數據", color=ft.Colors.RED)
        
        if self.page:
            self.page.update()
    
    def load_subscriptions(self):
        subscriptions = self.api_client.get_indicator_subscriptions(self.user_id)
        
        self.subscription_list.controls.clear()
        
        if not subscriptions:
            self.subscription_list.controls.append(
                ft.Container(
                    content=ft.Text("尚無訂閱", size=14, color=ft.Colors.GREY),
                    padding=20,
                )
            )
        else:
            for sub in subscriptions:
                enabled = sub.get("enabled", True)
                volume_check = sub.get("enableVolumeCheck", False)
                
                card = ft.Card(
                    content=ft.Container(
                        content=ft.Row([
                            ft.Column([
                                ft.Row([
                                    ft.Text(f"{sub.get('symbol', 'N/A')}", size=18, weight=ft.FontWeight.BOLD),
                                    ft.Container(
                                        content=ft.Text("啟用" if enabled else "停用", size=10),
                                        bgcolor=ft.Colors.GREEN_100 if enabled else ft.Colors.GREY_200,
                                        padding=3,
                                        border_radius=3,
                                    ),
                                ], spacing=10),
                                ft.Text(f"通知間隔: {sub.get('notifyIntervalMin', 60)} 分鐘", size=12, color=ft.Colors.GREY),
                                ft.Text(
                                    f"成交量判斷: {'開啟' if volume_check else '關閉'}", 
                                    size=12, 
                                    color=ft.Colors.BLUE if volume_check else ft.Colors.GREY
                                ),
                            ], spacing=3),
                            ft.Container(expand=True),
                            ft.Column([
                                ft.Switch(
                                    value=enabled,
                                    on_change=lambda e, sid=sub.get("subscriptionId"): self.toggle_subscription(sid),
                                ),
                                ft.IconButton(
                                    icon=ft.Icons.DELETE,
                                    icon_color=ft.Colors.RED,
                                    on_click=lambda e, sid=sub.get("subscriptionId"): self.delete_subscription(sid),
                                ),
                            ]),
                        ]),
                        padding=15,
                    ),
                )
                self.subscription_list.controls.append(card)
        
        if self.subscription_list.page:
            self.subscription_list.update()
    
    def toggle_subscription(self, subscription_id: str):
        result = self.api_client.toggle_indicator_subscription(subscription_id)
        if result:
            self.load_subscriptions()
            status = "啟用" if result.get("enabled") else "停用"
            self._show_snackbar(f"訂閱已{status}", ft.Colors.BLUE)
    
    def delete_subscription(self, subscription_id: str):
        if self.api_client.delete_indicator_subscription(subscription_id):
            self.load_subscriptions()
            self._show_snackbar("訂閱已刪除", ft.Colors.GREEN)
    
    def _show_snackbar(self, message: str, bgcolor):
        if self.page:
            snack_bar = ft.SnackBar(
                content=ft.Text(message, size=14),
                bgcolor=bgcolor,
            )
            self.page.overlay.append(snack_bar)
            snack_bar.open = True
            self.page.update()

