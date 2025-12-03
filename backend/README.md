# Backend - Go API Server

## 功能模組

1. **Price Fetcher Worker** - 定時抓取幣安價格存入 Redis
2. **Price API** - 提供價格查詢接口
3. **Alert Monitor** - 監控價格觸發警報推播

## 目錄結構

```
backend/
├── cmd/server/main.go           # 服務器入口
├── internal/
│   ├── api/                     # API 層
│   │   ├── handlers/            # HTTP 處理器
│   │   └── middleware/          # 中間件 (CORS, Gzip)
│   ├── worker/                  # 背景工作程序
│   ├── models/                  # 數據結構
│   ├── service/                 # 業務邏輯
│   └── repository/              # 數據訪問層
├── config/                      # 配置管理
└── go.mod                       # 依賴管理
```

## API 接口

- `GET /api/prices` - 獲取幣種價格列表
- `POST /api/alerts` - 創建價格警報
- `GET /api/alerts/:userId` - 查詢用戶警報
- `DELETE /api/alerts/:id` - 刪除警報

## 環境變數

```env
REDIS_URL=localhost:6379
FIREBASE_CREDENTIALS=path/to/serviceAccountKey.json
PORT=8080
```

