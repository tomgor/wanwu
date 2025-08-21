package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateAssistantCustom(ctx context.Context, assistantId uint32, customId string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 查重
		var count int64
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithCustomID(customId),
		).Apply(tx).Model(&model.AssistantCustom{}).Count(&count).Error; err != nil {
			return toErrStatus("assistant_custom_count", err.Error())
		}
		if count > 0 {
			return toErrStatus("assistant_custom_exist")
		}

		// 创建自定义工具
		if err := tx.Create(&model.AssistantCustom{
			AssistantId: assistantId,
			CustomId:    customId,
		}).Error; err != nil {
			return toErrStatus("assistant_custom_create", err.Error())
		}

		// 获取Assistant
		assistant := &model.Assistant{}
		if err := sqlopt.WithID(assistantId).Apply(tx).First(assistant).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}

		// 更新HasCustom字段
		if !assistant.HasCustom {
			if err := tx.Model(assistant).Update("has_custom", true).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) DeleteAssistantCustom(ctx context.Context, assistantId uint32, customId string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 1. 解析参数并构建查询条件
		cond := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithCustomID(customId),
		).Apply(tx)

		// 2. 检查数据是否存在
		var exists int64
		if err := cond.Model(&model.AssistantCustom{}).Count(&exists).Error; err != nil {
			return toErrStatus("assistant_custom_exist", err.Error())
		}
		if exists == 0 {
			return toErrStatus("assistant_custom_exist")
		}

		// 3. 执行删除操作
		if err := cond.Delete(&model.AssistantCustom{}).Error; err != nil {
			return toErrStatus("assistant_custom_delete", err.Error())
		}

		// 4. 检查剩余数量并更新状态
		var remaining int64
		if err := sqlopt.WithAssistantID(assistantId).Apply(tx).Model(&model.AssistantCustom{}).Count(&remaining).Error; err != nil {
			return toErrStatus("assistant_custom_count", err.Error())
		}

		if remaining == 0 {
			if err := sqlopt.WithID(assistantId).Apply(tx).Model(&model.Assistant{}).Update("has_custom", false).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) GetAssistantCustom(ctx context.Context, assistantId uint32, customId string) (*model.AssistantCustom, *err_code.Status) {
	custom := &model.AssistantCustom{}
	return custom, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(assistantId),
			sqlopt.WithCustomID(customId),
		).Apply(tx).First(custom).Error; err != nil {
			return toErrStatus("assistant_custom_get", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateAssistantCustom(ctx context.Context, custom *model.AssistantCustom) *err_code.Status {
	if custom.AssistantId == 0 {
		return toErrStatus("assistant_custom_update", "update assistant custom but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.SQLOptions(
			sqlopt.WithAssistantID(custom.AssistantId),
			sqlopt.WithCustomID(custom.CustomId),
		).Apply(tx).
			Model(&model.AssistantCustom{}).
			Update("enable", custom.Enable).
			Error; err != nil {
			return toErrStatus("assistant_custom_update", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistantCustomList(ctx context.Context, assistantId uint32) ([]*model.AssistantCustom, *err_code.Status) {
	var customList []*model.AssistantCustom
	return customList, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithAssistantID(assistantId).Apply(tx).Find(&customList).Error; err != nil {
			return toErrStatus("assistant_custom_list", err.Error())
		}
		return nil
	})
}
