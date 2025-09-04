package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
)

func (c *Client) CreateAssistantCustom(ctx context.Context, assistantId uint32, customId string, userId, orgID string) *err_code.Status {
	// 检查是否已存在
	var count int64
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantId),
		sqlopt.WithCustomID(customId),
	).Apply(c.db.WithContext(ctx)).Model(&model.AssistantCustom{}).
		Count(&count).Error; err != nil {
		return toErrStatus("assistant_custom_create", err.Error())
	}
	if count > 0 {
		return toErrStatus("assistant_custom_create", "custom already exists")
	}

	err := c.db.WithContext(ctx).Create(&model.AssistantCustom{
		AssistantId: assistantId,
		CustomId:    customId,
		Enable:      true, // 默认打开
		UserId:      userId,
		OrgId:       orgID,
	}).Error

	if err != nil {
		return toErrStatus("assistant_custom_create", err.Error())
	}
	return nil
}

func (c *Client) DeleteAssistantCustom(ctx context.Context, assistantId uint32, customId string) *err_code.Status {
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantId),
		sqlopt.WithCustomID(customId),
	).Apply(c.db.WithContext(ctx)).Delete(&model.AssistantCustom{}).Error; err != nil {
		return toErrStatus("assistant_custom_delete", err.Error())
	}
	return nil
}

func (c *Client) DeleteAssistantCustomByCustomToolId(ctx context.Context, customId string) *err_code.Status {
	if err := sqlopt.SQLOptions(
		sqlopt.WithCustomID(customId),
	).Apply(c.db.WithContext(ctx)).Delete(&model.AssistantCustom{}).Error; err != nil {
		return toErrStatus("assistant_custom_delete", err.Error())
	}
	return nil
}

func (c *Client) GetAssistantCustom(ctx context.Context, assistantId uint32, customId string) (*model.AssistantCustom, *err_code.Status) {
	custom := &model.AssistantCustom{}
	if err := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(assistantId),
		sqlopt.WithCustomID(customId),
	).Apply(c.db.WithContext(ctx)).First(custom).Error; err != nil {
		return nil, toErrStatus("assistant_custom_get", err.Error())
	}
	return custom, nil
}

func (c *Client) UpdateAssistantCustom(ctx context.Context, custom *model.AssistantCustom) *err_code.Status {
	// 更新
	result := sqlopt.SQLOptions(
		sqlopt.WithAssistantID(custom.AssistantId),
		sqlopt.WithCustomID(custom.CustomId),
	).Apply(c.db.WithContext(ctx)).
		Model(&model.AssistantCustom{}).
		Updates(map[string]interface{}{
			"enable": custom.Enable,
		})
	if result.Error != nil {
		return toErrStatus("assistant_custom_update", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return toErrStatus("assistant_custom_update", "custom not exists")
	}

	return nil
}

func (c *Client) GetAssistantCustomList(ctx context.Context, assistantId uint32) ([]*model.AssistantCustom, *err_code.Status) {
	var customList []*model.AssistantCustom
	if err := sqlopt.WithAssistantID(assistantId).Apply(c.db.WithContext(ctx)).Find(&customList).Error; err != nil {
		return nil, toErrStatus("assistant_custom_list", err.Error())
	}
	return customList, nil
}
