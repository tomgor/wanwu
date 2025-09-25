package response

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
)

type MCPSelect struct {
	MCPID       string `json:"mcpId"`       // mcpId
	MCPSquareID string `json:"mcpSquareId"` // 广场mcpId(非空表示来源于广场)
	Name        string `json:"name"`        // 名称
	UniqueId    string `json:"uniqueId"`    // 唯一标识
	Description string `json:"description"` // 描述
	ServerFrom  string `json:"serverFrom"`  // 来源
	ServerURL   string `json:"serverUrl"`   // sseUrl
}

type MCPToolList struct {
	Tools []*protocol.Tool `json:"tools"`
}

// MCPDetail MCP自定义详情
type MCPDetail struct {
	MCPInfo
	MCPSquareIntro
}

// MCPInfo MCP自定义信息
type MCPInfo struct {
	MCPID  string `json:"mcpId"`  // mcpId
	SSEURL string `json:"sseUrl"` // SSE URL
	MCPSquareInfo
}

// MCPSquareDetail MCP广场详情
type MCPSquareDetail struct {
	MCPSquareInfo
	MCPSquareIntro
	MCPTools
}

// MCPSquareInfo MCP广场信息
type MCPSquareInfo struct {
	MCPSquareID string         `json:"mcpSquareId"` // 广场mcpId(非空表示来源于广场)
	Avatar      request.Avatar `json:"avatar"`      // 图标
	Name        string         `json:"name"`        // 名称
	Desc        string         `json:"desc"`        // 描述
	From        string         `json:"from"`        // 来源
	Category    string         `json:"category"`    // 类型(data:数据,create:创作,search:搜索)
}

type MCPSquareIntro struct {
	Summary  string `json:"summary"`  // 使用概述
	Feature  string `json:"feature"`  // 特性说明
	Scenario string `json:"scenario"` // 应用场景
	Manual   string `json:"manual"`   // 使用说明
	Detail   string `json:"detail"`   // 详情
}

type MCPTools struct {
	SSEURL    string    `json:"sseUrl"`    // SSE URL
	Tools     []MCPTool `json:"tools"`     // 工具列表
	HasCustom bool      `json:"hasCustom"` // 是否已经发送到自定义
}

type MCPTool struct {
	Name        string             `json:"name"`        // 工具名
	Description string             `json:"description"` // 工具描述
	InputSchema MCPToolInputSchema `json:"inputSchema"` // 工具参数
}

type MCPToolInputSchema struct {
	Type       string                             `json:"type"`       // 固定值: object
	Properties map[string]MCPToolInputSchemaValue `json:"properties"` // 字段名 -> 字段信息
	Required   []string                           `json:"required"`   // 必填字段
}

type MCPToolInputSchemaValue struct {
	Type        string `json:"type"`        // 字段类型
	Description string `json:"description"` // 字段描述
}

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

// 精简OpenAPI结构体
type OpenAPI struct {
	OpenAPI string              `json:"openapi" yaml:"openapi"`
	Paths   map[string]PathItem `json:"paths" yaml:"paths"`
}

type PathItem struct {
	Get     *Operation `json:"get" yaml:"get"`
	Post    *Operation `json:"post" yaml:"post"`
	Put     *Operation `json:"put" yaml:"put"`
	Delete  *Operation `json:"delete" yaml:"delete"`
	Patch   *Operation `json:"patch" yaml:"patch"`
	Head    *Operation `json:"head" yaml:"head"`
	Options *Operation `json:"options" yaml:"options"`
}

type Operation struct {
	OperationID string `json:"operationId" yaml:"operationId"`
	Summary     string `json:"summary" yaml:"summary"`
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
