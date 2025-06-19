package request

type MCPCreate struct {
	MCPSquareID string `json:"mcpSquareId"`                // 广场mcpId(非空表示来源于广场)
	Name        string `json:"name" validate:"required"`   // 名称
	Desc        string `json:"desc" validate:"required"`   // 描述
	From        string `json:"from" validate:"required"`   // 来源
	SSEURL      string `json:"sseUrl" validate:"required"` // SSE URL
}
