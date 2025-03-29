package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func LineLoginURLHandler(c *gin.Context) {
	redirectURI := os.Getenv("DOMAIN") + "LineLogin"
	fmt.Println(redirectURI)
	state := "kirmcczgswokt024kqye0nx19n30o8nv"
	nonce := "rxz3j4i672bgqtxyu999hu4wkjc28de1"
	clientID := os.Getenv("LINE_client_id")

	// 轉換 redirect_uri 成 URL 編碼格式
	encodedRedirectURI := url.QueryEscape(redirectURI)

	// 建立完整的授權 URL
	authURL := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize"+
		"?response_type=code"+
		"&client_id=%s"+
		"&redirect_uri=%s"+
		"&state=%s"+
		"&bot_prompt=%s"+
		"&scope=profile%%20openid"+
		"&nonce=%s",
		clientID, encodedRedirectURI, state, "normal", nonce)

	c.IndentedJSON(http.StatusOK, gin.H{"url": authURL})
}

func SendMessage(ctx *gin.Context) {
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
	text := ctx.PostForm("text")
	userID := ctx.PostForm("userID")

	fmt.Println(text)
	fmt.Println(userID)
	// 建立訊息結構
	data := LineMessage{
		To: userID,
		Messages: []MessageBody{
			{Type: "text", Text: text},
		},
	}

	// 轉換為 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("❌ JSON Marshal Error:", err)
		return
	}

	// 建立 HTTP 請求
	req, err := http.NewRequest("POST", lineAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Request Creation Error:", err)
		return
	}

	// 設定 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request Error:", err)
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Message Sent:", string(body))
	} else {
		fmt.Printf("❌ LINE API Error: %d - %s\n", resp.StatusCode, string(body))
	}
}

func LineAuthHandler(ctx *gin.Context) {
	clientID := os.Getenv("LINE_client_id")
	clientSecret := os.Getenv("LINE_client_secret")
	// 從請求 body 讀取 code 和 state
	code := ctx.PostForm("code")
	//state := ctx.PostForm("state")

	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	// 設定 LINE Token 交換的參數
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("DOMAIN")+"LineLogin")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// 建立 HTTP POST 請求
	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 設定 Header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 送出請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Request failed"})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, _ := ioutil.ReadAll(resp.Body)

	// 解析並回傳 Token 給前端
	ctx.JSON(http.StatusOK, gin.H{"response": string(body)})
}
