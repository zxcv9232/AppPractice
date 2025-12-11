package indicators

// VolumeResult 成交量分析結果
type VolumeResult struct {
	CurrentVolume float64 // 當前 1 分 K 成交量
	AvgVolume     float64 // 近 N 根 K 線平均成交量
	VolumeRatio   float64 // 當前/平均 比值
}

// VolumeCheckMode 成交量檢查模式
type VolumeCheckMode string

const (
	VolumeCheckModeFixed      VolumeCheckMode = "fixed"      // 固定值模式
	VolumeCheckModeMultiplier VolumeCheckMode = "multiplier" // 倍數模式
)

// VolumeConfig 成交量檢查配置
type VolumeConfig struct {
	Enabled    bool            // 是否啟用成交量檢查
	Mode       VolumeCheckMode // 檢查模式
	FixedValue float64         // 固定值模式：閾值
	Multiplier float64         // 倍數模式：N 倍
	AvgPeriod  int             // 均量計算週期（幾根 K 線）
}

// DefaultVolumeConfig 返回預設配置
func DefaultVolumeConfig() VolumeConfig {
	return VolumeConfig{
		Enabled:    false,
		Mode:       VolumeCheckModeMultiplier,
		FixedValue: 0,
		Multiplier: 5.0,
		AvgPeriod:  60,
	}
}

// CalculateVolumeStats 計算成交量統計
// volumes: 成交量切片，最新的在最後
// avgPeriod: 計算平均的週期數
func CalculateVolumeStats(volumes []float64, avgPeriod int) VolumeResult {
	if len(volumes) == 0 {
		return VolumeResult{}
	}

	// 當前成交量（最後一根）
	currentVolume := volumes[len(volumes)-1]

	// 計算平均成交量（不含當前這根）
	avgVolume := 0.0
	count := 0

	// 從倒數第二根開始往前取 avgPeriod 根
	startIdx := len(volumes) - 1 - avgPeriod
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := len(volumes) - 1 // 不含最後一根

	for i := startIdx; i < endIdx; i++ {
		avgVolume += volumes[i]
		count++
	}

	if count > 0 {
		avgVolume = avgVolume / float64(count)
	}

	// 計算比值
	volumeRatio := 0.0
	if avgVolume > 0 {
		volumeRatio = currentVolume / avgVolume
	}

	return VolumeResult{
		CurrentVolume: currentVolume,
		AvgVolume:     avgVolume,
		VolumeRatio:   volumeRatio,
	}
}

// CheckVolumeCondition 檢查成交量是否滿足條件
func CheckVolumeCondition(result VolumeResult, config VolumeConfig) bool {
	if !config.Enabled {
		return true // 未啟用時直接通過
	}

	switch config.Mode {
	case VolumeCheckModeFixed:
		return result.CurrentVolume >= config.FixedValue
	case VolumeCheckModeMultiplier:
		return result.CurrentVolume >= (result.AvgVolume * config.Multiplier)
	default:
		return true
	}
}
