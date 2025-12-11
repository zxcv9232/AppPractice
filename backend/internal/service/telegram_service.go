package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cryptowatch/internal/models"

	"github.com/rs/zerolog/log"
)

// TelegramService Telegram Bot 通知服務
type TelegramService struct {
	botToken   string
	testMode   bool   // 測試模式：只 Log 不發送
	myChatID   string // 你自己的 Chat ID（測試用）
	apiBaseURL string
}

// TelegramMessage Telegram 發送訊息結構
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// TelegramResponse Telegram API 回應結構
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

// NewTelegramService 創建 Telegram 通知服務
func NewTelegramService(botToken string, testMode bool, myChatID string) *TelegramService {
	service := &TelegramService{
		botToken:   botToken,
		testMode:   testMode,
		myChatID:   myChatID,
		apiBaseURL: "https://api.telegram.org",
	}

	if testMode {
		log.Info().Msg("TelegramService running in test mode (log only)")
	} else if botToken == "" {
		log.Warn().Msg("TelegramService: No bot token provided, notifications will be logged only")
		service.testMode = true
	} else {
		log.Info().Msg("TelegramService initialized")
	}

	return service
}

// SendAlert 發送警報訊息
func (s *TelegramService) SendAlert(chatID string, payload models.AlertPayload) error {
	// 格式化訊息
	message := s.formatAlertMessage(payload)

	return s.sendMessage(chatID, message)
}

// SendMessage 發送一般訊息
func (s *TelegramService) sendMessage(chatID string, text string) error {
	// 測試模式
	if s.testMode {
		chatIDPreview := chatID
		if len(chatID) > 10 {
			chatIDPreview = chatID[:10] + "..."
		}
		log.Info().
			Str("chatId", chatIDPreview).
			Str("message", text).
			Msg("[TEST MODE] Would send Telegram message")
		return nil
	}

	// 沒有 Bot Token
	if s.botToken == "" {
		log.Warn().Msg("Telegram bot token not configured, skipping notification")
		return nil
	}

	// 構建請求
	url := fmt.Sprintf("%s/bot%s/sendMessage", s.apiBaseURL, s.botToken)

	msg := TelegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// 發送請求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Str("chatId", chatID).Msg("Failed to send Telegram message")
		return err
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(body, &telegramResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !telegramResp.OK {
		log.Error().
			Str("chatId", chatID).
			Str("error", telegramResp.Description).
			Msg("Telegram API error")
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	log.Info().
		Str("chatId", chatID).
		Msg("Telegram message sent successfully")

	return nil
}

// formatAlertMessage 格式化警報訊息
func (s *TelegramService) formatAlertMessage(payload models.AlertPayload) string {
	// 使用 HTML 格式
	message := fmt.Sprintf(
		"<b>%s</b>\n\n%s\n\n"+
			"📊 <b>詳細資訊</b>\n"+
			"├ 幣種: <code>%s</code>\n"+
			"├ 當前價格: <code>%.2f</code>\n"+
			"├ 上軌: <code>%.2f</code>\n"+
			"└ 下軌: <code>%.2f</code>",
		payload.Title,
		payload.Body,
		payload.Symbol,
		payload.CurrentPrice,
		payload.UpperBand,
		payload.LowerBand,
	)

	return message
}

// SendToMyself 發送到自己的 Chat ID（測試用）
func (s *TelegramService) SendToMyself(payload models.AlertPayload) error {
	if s.myChatID == "" {
		log.Warn().Msg("No personal Telegram chat ID configured")
		return nil
	}
	return s.SendAlert(s.myChatID, payload)
}

// IsEnabled 檢查通知服務是否可用
func (s *TelegramService) IsEnabled() bool {
	return s.botToken != "" || s.testMode
}

// SendRawMessage 發送原始文字訊息（不格式化）
func (s *TelegramService) SendRawMessage(chatID string, text string) error {
	return s.sendMessage(chatID, text)
}

