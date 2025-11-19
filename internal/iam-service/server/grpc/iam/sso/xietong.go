package sso

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/config"
)

const (
	API_ACCESS_TOKEN = "/api/portal-app/identity/getAccessToken"
	API_PARSE_TOKEN  = "/api/portal-app/identity/parseUserToken"
)

type ApiResponse struct {
	Code interface{}     `json:"code"` // 可能是 string "0" 或 int 0
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"` // 使用 RawMessage 延迟解析 Data 字段
}
type AccessTokenData struct {
	AccessToken string `json:"accessToken"`
}

// apiTransport 用于在请求中自动添加 clientId 和 Authorization 头部
type apiTransport struct {
	Transport   http.RoundTripper
	ClientID    string
	AccessToken string
}
type ParseTokenData struct {
	LoginName string `json:"loginName"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Org       struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"org"`
}

func FetchXieTongUserInfo(token string, ssoConfig *config.SSOConfig) (*model.User, error) {
	requestBody := map[string]any{"token": token}

	dataBuffer, err := apiWithToken(ssoConfig, API_PARSE_TOKEN, requestBody)
	if err != nil {
		return nil, err
	}

	// 解析 Data 字段为 UserInfo 结构体
	var data ParseTokenData
	if err := json.Unmarshal(*dataBuffer, &data); err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %w", err)
	}

	user := &model.User{
		Company:   data.Org.Name,
		Name:      data.LoginName,
		Nick:      data.Name,
		Phone:     data.Mobile,
		Email:     data.Email,
		CreatorID: 1,
		Remark:    "sso_xietong",
		Status:    true,
		IsAdmin:   false,
	}

	return user, errors.New("not implemented")
}

// apiWithToken 构造带有 clientId 和 Authorization 头的 http.Client
func apiWithToken(ssoConfig *config.SSOConfig, api string, requestBody map[string]any) (*[]byte, error) {

	accessToken, err := getAccessToken(ssoConfig)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s", ssoConfig.OtherProperties["base_url"].(string), api)
	log.Printf("请求地址：%s", url)

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("无法编码请求体为 JSON: %w", err)
	}
	log.Printf("请求体：%s", string(bodyJSON))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("无法创建 HTTP 请求: %w", err)
	}

	req.Header.Set("Content-Type", "application/json") // 声明请求体是 JSON 格式
	req.Header.Set("clientId", ssoConfig.OtherProperties["app_id"].(string))
	req.Header.Set("Authorization", accessToken)

	client, err := getHttpClient()
	if err != nil {
		return nil, fmt.Errorf("获取client失败: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}

	defer resp.Body.Close()

	// 3. 验证响应并解析 Data
	dataBuffer, err := validateResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("API响应验证失败: %w", err)
	}

	// Go 的 http.Client 需要通过请求对象设置 URL 和 headers
	// 这里返回一个带有自定义 Transport 的 Client 和目标 URL
	// 在 FetchUserInfo 中使用 client.Post(url, ...)
	return &dataBuffer, nil
}

// getAccessToken 获取 Access Token，使用缓存逻辑
func getAccessToken(ssoConfig *config.SSOConfig) (string, error) {

	// 1. 构建请求体
	requestBody := map[string]string{
		"grantType":    "client_credentials",
		"clientId":     ssoConfig.OtherProperties["app_id"].(string),
		"clientSecret": ssoConfig.OtherProperties["secret_key"].(string),
	}
	bodyJSON, _ := json.Marshal(requestBody)

	// 2. 发送请求 (POST /api/portal-app/identity/getAccessToken)
	url := fmt.Sprintf("%s%s", ssoConfig.OtherProperties["base_url"].(string), API_ACCESS_TOKEN)
	log.Printf("请求地址%s", url)
	log.Printf("请求体：%s", string(bodyJSON))

	client, err := getHttpClient()
	if err != nil {
		return "", fmt.Errorf("获取client失败: %w", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", fmt.Errorf("请求获取Token API失败: %w", err)
	}
	defer resp.Body.Close()

	// 3. 验证响应并解析 Data
	dataBuffer, err := validateResponse(resp)
	if err != nil {
		return "", fmt.Errorf("API响应验证失败: %w", err)
	}

	var data AccessTokenData
	if err := json.Unmarshal(dataBuffer, &data); err != nil {
		return "", fmt.Errorf("解析 Access Token 失败: %w", err)
	}

	return data.AccessToken, nil
}

func getHttpClient() (*http.Client, error) {
	tr := &http.Transport{
		// 配置 TLS 结构
		TLSClientConfig: &tls.Config{
			// 这将告诉客户端跳过验证服务器提供的证书
			InsecureSkipVerify: true,
		},
	}

	// 2. 创建一个使用自定义 Transport 的 HTTP 客户端
	client := &http.Client{
		Transport: tr,
		// 您可能还需要设置超时时间
		Timeout: 10 * time.Second,
	}
	return client, nil
}

// validateResponse 验证 HTTP 状态码和 JSON 响应体中的业务 Code
func validateResponse(resp *http.Response) ([]byte, error) {
	rspBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	log.Printf("返回：%s", string(rspBody))

	// 验证 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("返回码不对: %d", resp.StatusCode)
	}

	// 验证 JSON 格式
	var apiResp ApiResponse
	if err := json.Unmarshal(rspBody, &apiResp); err != nil {
		return nil, errors.New("返回数据非 JSON 格式或解析失败")
	}

	// 验证业务 Code
	msg := apiResp.Msg
	codeValid := false
	switch v := apiResp.Code.(type) {
	case string:
		if v == "0" {
			codeValid = true
		}
	case float64: // JSON numbers are parsed as float64 in Go's interface{}
		if int(v) == 0 {
			codeValid = true
		}
	case int:
		if v == 0 {
			codeValid = true
		}
	default:
		return nil, fmt.Errorf("未知的 code 类型: %T", apiResp.Code)
	}

	if !codeValid {
		return nil, fmt.Errorf("业务 Code 不为 0, 错误信息: %s", msg)
	}

	return apiResp.Data, nil
}
