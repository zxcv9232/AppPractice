import flet as ft

class ChartPage:
    def __init__(self, symbol: str, on_back):
        self.symbol = symbol
        self.on_back = on_back
    
    def build(self) -> ft.Container:
        tradingview_url = f"https://www.tradingview.com/chart/?symbol=BINANCE:{self.symbol}USDT"
        
        return ft.Container(
            content=ft.Column([
                ft.Container(
                    content=ft.Row([
                        ft.IconButton(
                            icon=ft.Icons.ARROW_BACK,
                            on_click=lambda e: self.on_back(),
                            tooltip="返回首頁"
                        ),
                        ft.Text(f"{self.symbol} K線圖", size=20, weight=ft.FontWeight.BOLD),
                    ]),
                    padding=10,
                ),
                ft.Container(
                    content=ft.Column([
                        ft.Text("📊 K 線圖功能", size=18, weight=ft.FontWeight.BOLD),
                        ft.Divider(),
                        ft.Text(f"幣種: {self.symbol}USDT", size=16),
                        ft.Text(f"TradingView 圖表網址:", size=14, color="grey"),
                        ft.Text(tradingview_url, size=12, selectable=True),
                        ft.Divider(),
                        ft.Text("💡 提示：", size=14, weight=ft.FontWeight.BOLD),
                        ft.Text("• 目前使用 TradingView 網址", size=12),
                        ft.Text("• WebView 需在打包成 App 後使用", size=12),
                        ft.Text("• 桌面版本可以複製網址到瀏覽器查看", size=12),
                        ft.Container(height=20),
                        ft.ElevatedButton(
                            text="在瀏覽器中打開",
                            icon=ft.Icons.OPEN_IN_BROWSER,
                            on_click=lambda e: self.open_in_browser(tradingview_url)
                        ),
                    ], 
                    horizontal_alignment=ft.CrossAxisAlignment.START,
                    spacing=10),
                    expand=True,
                    padding=20,
                ),
            ]),
        )
    
    def open_in_browser(self, url: str):
        import webbrowser
        webbrowser.open(url)

