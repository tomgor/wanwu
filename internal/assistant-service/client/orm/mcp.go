package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
)

func (c *Client) CreateAssistantMCP(ctx context.Context, assistantId uint32, mcpId string, userId, orgID string) *err_code.Status {
	if err := c.db.Create(&model.AssistantMCP{
		AssistantId: assistantId,
		MCPId:       mcpId,
		Enable:      true, // 默认开
		UserId:      userId,
		OrgId:       orgID,
	}).Error; err != nil {
		return toErrStatus("assistant_mcp_create", err.Error())
	}
	return nil
}

func (c *Client) DeleteAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) *err_code.Status {
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantId),
		sqlopt.WithMCPID(mcpId),
	).Apply(c.db).Delete(&model.AssistantMCP{}).Error; err != nil {
		return toErrStatus("assistant_mcp_delete", err.Error())
	}
	return nil
}

func (c *Client) GetAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) (*model.AssistantMCP, *err_code.Status) {
	mcp := &model.AssistantMCP{}
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantId),
		sqlopt.WithMCPID(mcpId),
	).Apply(c.db).First(mcp).Error; err != nil {
		return nil, toErrStatus("assistant_mcp_get", err.Error())
	}
	return mcp, nil
}

func (c *Client) GetAssistantMCPList(ctx context.Context, assistantId uint32) ([]*model.AssistantMCP, *err_code.Status) {
	var mcpList []*model.AssistantMCP
	if err := sqlopt.WithAssistantID(assistantId).Apply(c.db).Find(&mcpList).Error; err != nil {
		return nil, toErrStatus("assistant_mcp_list", err.Error())
	}
	return mcpList, nil
}

func (c *Client) UpdateAssistantMCP(ctx context.Context, mcp *model.AssistantMCP) *err_code.Status {
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(mcp.AssistantId),
		sqlopt.WithMCPID(mcp.MCPId),
	).Apply(c.db).Model(&model.AssistantMCP{}).Updates(map[string]interface{}{
		"enable": mcp.Enable,
	}).Error; err != nil {
		return toErrStatus("assistant_mcp_update", err.Error())
	}
	return nil
}
