package orm

import (
	"context"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
	"time"
)

func (c *Client) CreateAssistantMCP(ctx context.Context, mcp *model.AssistantMCP) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 查重
		var count int64
		if err := tx.Model(&model.AssistantMCP{}).Where("assistant_id = ? AND mcp_id = ?", mcp.AssistantId, mcp.MCPId).Count(&count).Error; err != nil {
			return toErrStatus("assistant_mcp_count", err.Error())
		}
		if count > 0 {
			return toErrStatus("assistant_mcp_exist", "MCP已存在")
		}

		// 创建MCP
		if err := tx.Create(mcp).Error; err != nil {
			return toErrStatus("assistant_mcp_create", err.Error())
		}

		// 获取Assistant
		assistant := &model.Assistant{}
		if err := sqlopt.WithID(mcp.AssistantId).Apply(tx).First(assistant).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}

		// 更新HasMCP字段
		if !assistant.HasMCP {
			assistant.HasMCP = true
			if err := tx.Model(assistant).Updates(map[string]interface{}{
				"has_mcp": assistant.HasMCP,
			}).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) DeleteAssistantMCP(ctx context.Context, id uint32) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 获取MCP以获得AssistantId
		mcp := &model.AssistantMCP{}
		if err := sqlopt.WithID(id).Apply(tx).First(mcp).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}
		assistantId := mcp.AssistantId

		// 删除MCP
		if err := sqlopt.WithID(id).Apply(tx).Delete(&model.AssistantMCP{}).Error; err != nil {
			return toErrStatus("assistant_mcp_delete", err.Error())
		}

		// 查询是否还有其他mcp
		var count int64
		if err := tx.Where("assistant_id = ?", assistantId).Model(&model.AssistantMCP{}).Count(&count).Error; err != nil {
			return toErrStatus("assistant_mcp_count", err.Error())
		}

		// 如果没有其他mcp，更新Assistant的HasMCP字段
		if count == 0 {
			assistant := &model.Assistant{}
			if err := sqlopt.WithID(assistantId).Apply(tx).First(assistant).Error; err != nil {
				return toErrStatus("assistant_get", err.Error())
			}
			assistant.HasMCP = false
			if err := tx.Model(assistant).Updates(map[string]interface{}{
				"has_mcp": assistant.HasMCP,
			}).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) GetAssistantMCP(ctx context.Context, query map[string]interface{}) (*model.AssistantMCP, *err_code.Status) {
	var mcp *model.AssistantMCP
	return mcp, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		mcp = &model.AssistantMCP{}

		for key, value := range query {
			tx = tx.Where(key, value)
		}

		if err := tx.First(mcp).Error; err != nil {
			return toErrStatus("assistant_mcp_get", err.Error())
		}

		return nil
	})
}

func (c *Client) GetAssistantMCPList(ctx context.Context, query map[string]interface{}) ([]*model.AssistantMCP, *err_code.Status) {
	var mcpList []*model.AssistantMCP
	return mcpList, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		mcpList = []*model.AssistantMCP{}

		for key, value := range query {
			tx = tx.Where(key, value)
		}

		if err := tx.Find(&mcpList).Error; err != nil {
			return toErrStatus("assistant_mcp_list", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateAssistantMCP(ctx context.Context, mcp *model.AssistantMCP) *err_code.Status {
	if mcp.ID == 0 {
		return toErrStatus("assistant_mcp_update", "update assistant mcp but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Model(mcp).Updates(map[string]interface{}{
			"enable":     mcp.Enable,
			"updated_at": time.Now().UnixMilli(),
		}).Error; err != nil {
			return toErrStatus("assistant_mcp_update", err.Error())
		}
		return nil
	})
}
