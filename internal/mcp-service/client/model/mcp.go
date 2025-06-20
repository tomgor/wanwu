package model

type MCPModel struct {
	ID          uint32 `gorm:"primary_key;auto_increment;not null;"`
	SseUrl      string `json:"serverUrl" gorm:"column:sse_url;type:text;comment:服务地址"`
	McpSquareId string `json:"mcpSquareId" gorm:"column:mcp_square_id;type:varchar(100);comment:广场mcpId"`
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null;default:'';comment:服务名称"`
	Desc        string `json:"desc" gorm:"column:desc;type:text;comment:功能描述"`
	From        string `json:"from" gorm:"column:from;type:text;comment:服务来源"`
	PublicModel
}

type PublicModel struct {
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli;index:create_at;column:created_at;type:bigint(20);comment:创建时间"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoCreateTime:milli;column:updated_at;type:bigint(20);comment:更新时间"`
	OrgID     string `gorm:"index:org_id;column:org_id;type:varchar(255);comment:组织ID" json:"orgId"`
	UserID    string `gorm:"index:user_id;column:user_id;type:varchar(255);comment:用户ID" json:"userId"`
}

func (MCPModel) TableName() string {
	return "mcp_info"
}
