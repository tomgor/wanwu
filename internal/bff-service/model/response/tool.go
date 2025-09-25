package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type CustomToolApiAuthWebRequest struct {
	Type             string `json:"type"`             // 认证类型: None 或 APIKey
	APIKey           string `json:"apiKey"`           // apiKey
	CustomHeaderName string `json:"customHeaderName"` // 自定义头名
	AuthType         string `json:"authType"`         // Auth类型
}

type CustomToolApiResponse struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

type CustomToolDetail struct {
	CustomToolId  string                      `json:"customToolId"`  // 自定义工具id
	Name          string                      `json:"name"`          // 名称
	Description   string                      `json:"description"`   // 描述
	Schema        string                      `json:"schema"`        // schema
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth"`       // apiAuth
	ApiList       []CustomToolApiResponse     `json:"apiList"`       // api列表
	PrivacyPolicy string                      `json:"privacyPolicy"` // 隐私政策
	ToolSquareID  string                      `json:"toolSquareId"`  // 广场mcpId(非空表示来源于广场)
}

type CustomToolCell struct {
	CustomToolId string `json:"customToolId"` // 自定义工具id
	Name         string `json:"name"`         // 名称
	Description  string `json:"description"`  // 描述
}

type CustomToolSelect struct {
	UniqueId     string `json:"uniqueId"`     // 统一的id
	CustomToolId string `json:"customToolId"` // 自定义工具id
	Name         string `json:"name"`         // 名称
	Description  string `json:"description"`  // 描述
}

type ToolSquareInfo struct {
	ToolSquareID string         `json:"toolSquareId"` // 广场mcpId(非空表示来源于广场)
	Avatar       request.Avatar `json:"avatar"`       // 图标
	Name         string         `json:"name"`         // 名称
	Desc         string         `json:"desc"`         // 描述
	Tags         []string       `json:"tags"`         // 标签
}

type BuiltInTools struct {
	NeedApiKeyInput bool      `json:"needApiKeyInput"` // 是否需要apiKey输入
	APIKey          string    `json:"apiKey"`          // apiKey
	Tools           []MCPTool `json:"tools"`           // 工具列表
	Detail          string    `json:"detail"`          // 详细描述
	ActionSum       int64     `json:"actionSum"`       // action总数
}

type ToolSquareDetail struct {
	ToolSquareInfo
	BuiltInTools
}
