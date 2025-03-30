package service

import (
	"bots/shop/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// 用戶行為服務層

// 購買
func Buy(userID string, commodity_id, commodity_spec_id, num uint) error {
	// 先扣除數量
	commodity_spec := models.GetCommoditySpec(commodity_spec_id)

	if _, err := models.UpdateCommoditySpecification(commodity_id, commodity_spec_id, models.UpdateCommoditySpecRequest{
		Stock:      func(v uint) *uint { return &v }(commodity_spec.Stock - num),
		Price:      &commodity_spec.Price,
		PictureURL: commodity_spec.PictureUrl,
	}); err != nil {
		return err
	}

	// 若成功發送LineMessage
	message := fmt.Sprintf("您已購買商品: %s, 規格: %s %s, 數量: %d, 總價格為: %.2f", commodity_spec.Commodity.CommodityName, commodity_spec.SpecValue1.SpecValue, commodity_spec.SpecValue2.SpecValue, num, (commodity_spec.Price * float64(num)))
	if err := SendMessageToUser(message, userID); err != nil {
		return err
	}
	return nil
}

// 發送訊息給用戶
func SendMessageToUser(message string, userID string) error {
	token := os.Getenv("LINE_ACCESS_TOKEN")

	type MessageBody struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	type LineMessage struct {
		To       string        `json:"to"`
		Messages []MessageBody `json:"messages"`
	}
	lineAPIURL := "https://api.line.me/v2/bot/message/push"

	// 建立訊息結構
	data := LineMessage{
		To: userID,
		Messages: []MessageBody{
			{Type: "text", Text: message},
		},
	}

	// 轉換為 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("❌ JSON Marshal Error:", err)
		return err
	}

	// 建立 HTTP 請求
	req, err := http.NewRequest("POST", lineAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Request Creation Error:", err)
		return err
	}

	// 設定 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
