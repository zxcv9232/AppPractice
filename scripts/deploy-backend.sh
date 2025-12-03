#!/bin/bash

echo "🚀 準備部署後端到 Render..."

cd "$(dirname "$0")/.."

if [ ! -d ".git" ]; then
    echo "📦 初始化 Git repository..."
    git init
    git add .
    git commit -m "Initial commit for deployment"
fi

echo ""
echo "✅ 後端代碼已準備就緒！"
echo ""
echo "接下來的步驟："
echo "1. 在 GitHub 上創建新的 repository"
echo "2. 執行以下命令推送代碼："
echo ""
echo "   git remote add origin YOUR_GITHUB_REPO_URL"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "3. 前往 https://render.com 註冊帳號"
echo "4. 創建 Redis 服務（免費版）"
echo "5. 創建 Web Service 並連接您的 GitHub repo"
echo "6. 設定環境變數："
echo "   - PORT=8080"
echo "   - REDIS_URL=(從 Redis 服務複製)"
echo ""
echo "7. 等待部署完成，複製您的服務 URL"
echo "8. 更新前端 api_client.py 中的 base_url"
echo ""

