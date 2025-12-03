#!/bin/bash

echo "🌐 啟動虛擬幣看盤 Web 版本"
echo ""

cd "$(dirname "$0")/../frontend"

if [ ! -f "requirements.txt" ]; then
    echo "❌ 找不到 requirements.txt"
    exit 1
fi

if ! python -c "import flet" 2>/dev/null; then
    echo "📦 安裝依賴套件..."
    pip install -r requirements.txt
fi

echo "🚀 啟動 Web 服務器..."
echo ""
echo "✅ 服務器啟動後，請在瀏覽器中訪問:"
echo "   👉 http://localhost:8080"
echo ""
echo "提示："
echo "  - 確保後端 API 正在運行 (go run cmd/server/main.go)"
echo "  - 按 Ctrl+C 停止服務器"
echo ""

python src/main.py --web

