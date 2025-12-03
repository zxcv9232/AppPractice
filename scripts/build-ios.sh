#!/bin/bash

echo "📱 開始打包 iOS App..."

cd "$(dirname "$0")/../frontend"

if ! command -v flet &> /dev/null; then
    echo "❌ Flet 未安裝，正在安裝..."
    pip install flet
fi

echo ""
echo "請確認以下資訊："
echo "- Bundle ID: com.yourname.cryptowatch"
echo "- App 名稱: 虛擬幣看盤"
echo "- 版本號: 1.0.0"
echo ""
read -p "按 Enter 繼續打包，或 Ctrl+C 取消..."

echo ""
echo "🔨 打包中，這可能需要幾分鐘..."
flet build ipa

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 打包成功！"
    echo ""
    echo "IPA 檔案位置："
    find . -name "*.ipa" -type f
    echo ""
    echo "接下來的步驟："
    echo "1. 下載並打開 Transporter App"
    echo "2. 拖曳 .ipa 檔案到 Transporter"
    echo "3. 點擊 Deliver 上傳到 App Store Connect"
    echo "4. 前往 https://appstoreconnect.apple.com"
    echo "5. 在 TestFlight 中測試"
    echo "6. 提交 App Store 審核"
else
    echo ""
    echo "❌ 打包失敗，請檢查錯誤訊息"
    echo ""
    echo "常見問題："
    echo "- 是否已安裝 Xcode？"
    echo "- 是否已登入 Apple Developer 帳號？"
    echo "- Bundle ID 是否已在 Apple Developer 網站註冊？"
fi

