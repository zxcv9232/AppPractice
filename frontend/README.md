# Frontend - Python Flet App

## 功能頁面

1. **首頁** - 虛擬貨幣價格列表 (下拉刷新)
2. **K線圖頁** - TradingView 圖表 (WebView)
3. **警報設置頁** - 創建/管理價格與成交量警報
   - 價格警報：監控價格突破或跌破
   - 成交量警報：監控累積成交量達標

## 目錄結構

```
frontend/
├── src/
│   ├── main.py                  # App 入口
│   ├── ui/                      # UI 組件
│   │   ├── home_page.py         # 首頁列表
│   │   ├── chart_page.py        # K線圖頁
│   │   └── alert_page.py        # 警報管理頁
│   ├── services/                # 服務層
│   │   ├── api_client.py        # 後端 API 調用
│   │   └── firebase_client.py   # Firebase 認證
│   └── utils/                   # 工具函數
│       └── helpers.py           # 輔助函數
├── assets/                      # 靜態資源 (圖標、圖片)
└── requirements.txt             # Python 依賴
```

## 依賴套件

- `flet` - UI 框架
- `requests` - HTTP 客戶端
- `firebase-admin` - Firebase SDK

## 啟動應用

```bash
python src/main.py
```

## 警報功能說明

### 價格警報
設定當價格達到目標時觸發通知：
- 選擇「價格警報」標籤
- 輸入幣種、目標價格和觸發條件
- 支援「突破」和「跌破」兩種條件

### 成交量警報 ⭐ 新功能
監控指定時間窗口內的累積成交量：
- 選擇「成交量警報」標籤
- 輸入幣種、目標成交量和時間窗口（分鐘）
- 當累積成交量達到目標時觸發通知

詳細說明請參考：[前端成交量警報功能文檔](../docs/frontend-volume-alerts.md)

