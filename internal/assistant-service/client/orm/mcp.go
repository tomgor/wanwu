package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 查重
		var count int64
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithMCPID(mcpId),
		).Apply(tx).Model(&model.AssistantMCP{}).Count(&count).Error; err != nil {
			return toErrStatus("assistant_mcp_count", err.Error())
		}
		if count > 0 {
			return toErrStatus("assistant_mcp_exist", "mcp already exist")
		}

		// 创建MCP
		if err := tx.Create(&model.AssistantMCP{
			AssistantId: assistantId,
			MCPId:       mcpId,
		}).Error; err != nil {
			return toErrStatus("assistant_mcp_create", err.Error())
		}

		// 获取Assistant
		assistant := &model.Assistant{}
		if err := sqlopt.WithID(assistantId).Apply(tx).First(assistant).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}

		// 更新HasMCP字段
		if !assistant.HasMCP {
			if err := tx.Model(assistant).Update("has_mcp", true).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) DeleteAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 删除MCP
		cond := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithMCPID(mcpId),
		).Apply(tx)

		var exists int64
		if err := cond.Model(&model.AssistantMCP{}).Count(&exists).Error; err != nil {
			return toErrStatus("assistant_mcp_delete", err.Error())
		}
		if exists == 0 {
			return toErrStatus("assistant_mcp_delete", "mcp not exist")
		}

		if err := cond.Delete(&model.AssistantMCP{}).Error; err != nil {
			return toErrStatus("assistant_mcp_delete", err.Error())
		}

		// 查询是否还有其他mcp
		var count int64
		if err := sqlopt.WithAssistantID(assistantId).Apply(tx).Model(&model.AssistantMCP{}).Count(&count).Error; err != nil {
			return toErrStatus("assistant_mcp_count", err.Error())
		}

		// 如果没有其他mcp，更新Assistant的HasMCP字段
		if count == 0 {
			if err := sqlopt.WithID(assistantId).Apply(tx).Model(&model.Assistant{}).Update("has_mcp", false).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) GetAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) (*model.AssistantMCP, *err_code.Status) {
	mcp := &model.AssistantMCP{}
	return mcp, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithMCPID(mcpId),
		).Apply(tx).First(mcp).Error; err != nil {
			return toErrStatus("assistant_mcp_get", err.Error())
		}

		return nil
	})
}

func (c *Client) GetAssistantMCPList(ctx context.Context, assistantId uint32) ([]*model.AssistantMCP, *err_code.Status) {
	var mcpList []*model.AssistantMCP
	return mcpList, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithAssistantID(assistantId).Apply(tx).Find(&mcpList).Error; err != nil {
			return toErrStatus("assistant_mcp_list", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateAssistantMCP(ctx context.Context, mcp *model.AssistantMCP) *err_code.Status {
	if mcp.AssistantId == 0 {
		return toErrStatus("assistant_mcp_update", "update assistant mcp but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(mcp.AssistantId),
			sqlopt.WithMCPID(mcp.MCPId),
		).Apply(tx).Model(&model.AssistantMCP{}).Update("enable", mcp.Enable).Error; err != nil {
			return toErrStatus("assistant_mcp_update", err.Error())
		}
		return nil
	})
}
