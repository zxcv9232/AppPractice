#!/bin/bash

API_URL="http://localhost:8080/api"

echo "=== 測試成交量監控功能 ==="
echo ""

echo "1. 創建成交量警報 (BTC 5分鐘內成交量達到 1,000,000)"
curl -X POST "$API_URL/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "test-user-1",
    "symbol": "BTC",
    "alertType": "volume",
    "targetVolume": 1000000,
    "timeWindow": 5
  }'
echo -e "\n"

echo "2. 創建成交量警報 (ETH 1分鐘內成交量達到 500,000)"
curl -X POST "$API_URL/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "test-user-1",
    "symbol": "ETH",
    "alertType": "volume",
    "targetVolume": 500000,
    "timeWindow": 1
  }'
echo -e "\n"

echo "3. 創建價格警報 (BTC 價格超過 50,000)"
curl -X POST "$API_URL/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "test-user-1",
    "symbol": "BTC",
    "alertType": "price",
    "targetPrice": 50000,
    "direction": "above"
  }'
echo -e "\n"

echo "4. 查詢用戶的所有警報"
curl -X GET "$API_URL/alerts/test-user-1"
echo -e "\n"

echo "5. 查詢當前價格（包含成交量資訊）"
curl -X GET "$API_URL/prices"
echo -e "\n"

echo "=== 測試完成 ==="

