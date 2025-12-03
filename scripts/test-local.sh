#!/bin/bash

echo "🧪 啟動本地測試環境..."
echo ""

check_redis() {
    if redis-cli ping > /dev/null 2>&1; then
        echo "✅ Redis 運行中"
        return 0
    else
        echo "❌ Redis 未運行"
        return 1
    fi
}

if ! check_redis; then
    echo ""
    echo "正在啟動 Redis..."
    if command -v brew &> /dev/null; then
        brew services start redis
        sleep 2
        if check_redis; then
            echo "✅ Redis 已啟動"
        else
            echo "❌ Redis 啟動失敗，請手動安裝: brew install redis"
            exit 1
        fi
    else
        echo "請安裝 Homebrew 或手動啟動 Redis"
        exit 1
    fi
fi

echo ""
echo "🚀 啟動後端服務..."
cd "$(dirname "$0")/../backend"

if [ ! -f "go.mod" ]; then
    echo "❌ 找不到 go.mod，請確認在正確的目錄"
    exit 1
fi

osascript -e 'tell app "Terminal"
    do script "cd '"$(pwd)"' && go mod download && go run cmd/server/main.go"
end tell'

echo "✅ 後端正在啟動... (在新終端視窗)"
echo ""
echo "等待 3 秒讓後端啟動..."
sleep 3

echo ""
echo "🧪 測試後端 API..."
if curl -s http://localhost:8080/api/prices > /dev/null; then
    echo "✅ 後端 API 運行正常"
else
    echo "⚠️  後端 API 尚未就緒，請稍候再試"
fi

echo ""
echo "📱 啟動前端..."
cd ../frontend

if [ ! -f "requirements.txt" ]; then
    echo "❌ 找不到 requirements.txt"
    exit 1
fi

osascript -e 'tell app "Terminal"
    do script "cd '"$(pwd)"' && python src/main.py"
end tell'

echo "✅ 前端正在啟動... (在新終端視窗)"
echo ""
echo "✨ 測試環境已啟動！"
echo ""
echo "您應該可以看到："
echo "1. 後端終端顯示 'Server starting on port 8080'"
echo "2. 前端視窗開啟虛擬幣看盤 App"
echo ""
echo "測試項目："
echo "- [ ] 首頁顯示價格列表"
echo "- [ ] 點擊幣種可查看 K 線圖"
echo "- [ ] 可以創建價格警報"
echo "- [ ] 警報列表正常顯示"

