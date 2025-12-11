package indicators

import (
	"fmt"
	"math"
)

// LRCResult 用來儲存線性回歸通道計算結果
type LRCResult struct {
	CenterLine float64 // 回歸中線 (LinReg)
	UpperBand  float64 // 上軌 (中線 + dev * deviation)
	LowerBand  float64 // 下軌 (中線 - dev * deviation)
	Slope      float64 // 斜率
	Deviation  float64 // 標準差
}

// LRCConfig LRC 計算配置
type LRCConfig struct {
	Length        int     // 回歸長度，預設 42
	DevMultiplier float64 // 標準差倍數，預設 2.0
}

// DefaultLRCConfig 返回預設配置
func DefaultLRCConfig() LRCConfig {
	return LRCConfig{
		Length:        42,
		DevMultiplier: 2.0,
	}
}

// CalculateLRC 計算線性回歸通道
// prices: 收盤價切片 (請確保長度至少等於 length)
// length: 回歸長度 (例如 42)
// devMultiplier: 標準差倍數 (例如 2.0)
func CalculateLRC(prices []float64, length int, devMultiplier float64) (LRCResult, error) {
	if len(prices) < length {
		return LRCResult{}, fmt.Errorf("數據長度不足，需要至少 %d 筆數據，目前只有 %d 筆", length, len(prices))
	}

	// 取出最後 length 筆數據進行計算
	// window[0] 是最舊的數據 (x=0)，window[length-1] 是最新的數據 (x=length-1)
	window := prices[len(prices)-length:]

	var sumX, sumY, sumXY, sumXX float64
	n := float64(length)

	// 1. 計算線性回歸所需的總和 (Least Squares)
	for i := 0; i < length; i++ {
		x := float64(i)
		y := window[i]

		sumX += x
		sumY += y
		sumXY += (x * y)
		sumXX += (x * x)
	}

	// 2. 計算斜率 (Slope) 和 截距 (Intercept)
	// 公式: Slope = (n*Σxy - Σx*Σy) / (n*Σx^2 - (Σx)^2)
	denominator := n*sumXX - sumX*sumX
	if denominator == 0 {
		return LRCResult{}, fmt.Errorf("無法計算斜率，分母為零")
	}
	slope := (n*sumXY - sumX*sumY) / denominator

	// 公式: Intercept = (Σy - slope*Σx) / n
	intercept := (sumY - slope*sumX) / n

	// 3. 計算當前 K 棒的回歸值 (Center Line)
	// 在這個 window 中，最新的 K 棒位置 x = length - 1
	currentX := float64(length - 1)
	centerLine := slope*currentX + intercept

	// 4. 計算標準差 (Deviation)
	// Pine Script 中的邏輯是計算「價格」與「回歸線上對應點」距離的平方和
	var sumResidualsSq float64
	for i := 0; i < length; i++ {
		x := float64(i)
		y := window[i]

		// 預測值 (回歸線上的點)
		yPred := slope*x + intercept

		// 殘差平方
		sumResidualsSq += math.Pow(y-yPred, 2)
	}

	// 標準差 = sqrt(殘差平方和 / n)
	deviation := math.Sqrt(sumResidualsSq / n)

	// 5. 計算上下軌
	upperBand := centerLine + (deviation * devMultiplier)
	lowerBand := centerLine - (deviation * devMultiplier)

	return LRCResult{
		CenterLine: centerLine,
		UpperBand:  upperBand,
		LowerBand:  lowerBand,
		Slope:      slope,
		Deviation:  deviation,
	}, nil
}

// CalculateLRCWithConfig 使用配置計算 LRC
func CalculateLRCWithConfig(prices []float64, config LRCConfig) (LRCResult, error) {
	return CalculateLRC(prices, config.Length, config.DevMultiplier)
}

