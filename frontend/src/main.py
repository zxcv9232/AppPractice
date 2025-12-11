import flet as ft
import sys
from services.api_client import APIClient
from services.firebase_client import FirebaseClient
from ui.home_page import HomePage
from ui.chart_page import ChartPage
from ui.alert_page import AlertPage
from ui.indicator_page import IndicatorPage

def main(page: ft.Page):
    page.title = "虛擬幣看盤 App"
    page.theme_mode = ft.ThemeMode.LIGHT
    page.window_width = 400
    page.window_height = 800
    
    api_client = APIClient()
    firebase_client = FirebaseClient()
    user_id = firebase_client.sign_in_anonymously()
    
    def show_home():
        home_page.stop_auto_refresh()
        page.controls.clear()
        page.add(home_page.build())
        page.add(navigation_bar)
        page.update()
        home_page.load_prices()
        home_page.start_auto_refresh()
    
    def show_chart(symbol: str):
        home_page.stop_auto_refresh()
        chart_page = ChartPage(symbol, on_back=show_home)
        page.controls.clear()
        page.add(chart_page.build())
        page.update()
    
    home_page = HomePage(api_client, show_chart)
    alert_page = AlertPage(api_client, user_id)
    indicator_page = IndicatorPage(api_client, user_id)
    
    alert_page.set_home_page(home_page)
    
    def on_navigation_change(e):
        selected_index = e.control.selected_index
        page.controls.clear()
        
        if selected_index == 0:
            home_page.stop_auto_refresh()
            page.add(home_page.build())
            page.add(navigation_bar)
            page.update()
            home_page.load_prices()
            home_page.start_auto_refresh()
        elif selected_index == 1:
            home_page.stop_auto_refresh()
            page.add(alert_page.build())
            page.add(navigation_bar)
            page.update()
            alert_page.load_alerts(check_triggered=False)
        elif selected_index == 2:
            home_page.stop_auto_refresh()
            indicator_page.set_page(page)
            page.add(indicator_page.build())
            page.add(navigation_bar)
            page.update()
            indicator_page.load_subscriptions()
    
    navigation_bar = ft.NavigationBar(
        destinations=[
            ft.NavigationBarDestination(icon=ft.Icons.HOME, label="首頁"),
            ft.NavigationBarDestination(icon=ft.Icons.NOTIFICATIONS, label="警報"),
            ft.NavigationBarDestination(icon=ft.Icons.SHOW_CHART, label="指標監控"),
        ],
        on_change=on_navigation_change,
    )
    
    page.add(home_page.build())
    page.add(navigation_bar)
    home_page.load_prices()
    home_page.start_auto_refresh()
    
    alert_page.start_monitoring(page)
    print("🌍 全局警報監控已啟動 - 在任何頁面都會顯示通知")

if __name__ == "__main__":
    if "--web" in sys.argv:
        print("🌐 啟動 Web 版本...")
        print("📍 訪問: http://localhost:8080")
        ft.app(target=main, view=ft.AppView.WEB_BROWSER, port=8080)
    else:
        print("📱 啟動桌面/App 版本...")
        ft.app(target=main)

