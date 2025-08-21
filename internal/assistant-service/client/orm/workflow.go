package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateAssistantWorkflow(ctx context.Context, workflow *model.AssistantWorkflow) *err_code.Status {
	if workflow.ID != 0 {
		return toErrStatus("assistant_workflow_create", "create assistant workflow but id not 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		// 创建Workflow
		if err := tx.Create(workflow).Error; err != nil {
			return toErrStatus("assistant_workflow_create", err.Error())
		}

		// 获取Assistant
		assistant := &model.Assistant{}
		if err := sqlopt.WithID(workflow.AssistantId).Apply(tx).First(assistant).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}

		// 更新HasWorkflow字段
		if !assistant.HasWorkflow {
			if err := tx.Model(assistant).Update("has_workflow", true).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}

func (c *Client) UpdateAssistantWorkflow(ctx context.Context, workflow *model.AssistantWorkflow) *err_code.Status {
	if workflow.ID == 0 {
		return toErrStatus("assistant_workflow_update", "update assistant workflow but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		cond := sqlopt.SQLOptions(
			sqlopt.WithID(workflow.AssistantId),
			sqlopt.WithWorkflowID(workflow.WorkflowId),
		).Apply(tx)
		if err := cond.Model(workflow).Update("enable", workflow.Enable).Error; err != nil {
			return toErrStatus("assistant_workflow_update", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistantWorkflow(ctx context.Context, assistantId uint32, workflowId string) (*model.AssistantWorkflow, *err_code.Status) {
	workflow := &model.AssistantWorkflow{}
	return workflow, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		cond := sqlopt.SQLOptions(
			sqlopt.WithID(assistantId),
			sqlopt.WithWorkflowID(workflowId),
		).Apply(tx)
		if err := cond.First(workflow).Error; err != nil {
			return toErrStatus("assistant_workflow_get", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistantWorkflowsByAssistantID(ctx context.Context, assistantId uint32) ([]*model.AssistantWorkflow, *err_code.Status) {
	var workflows []*model.AssistantWorkflow
	return workflows, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithAssistantID(assistantId).Apply(tx).Find(&workflows).Error; err != nil {
			return toErrStatus("assistant_workflows_get_by_assistant_id", err.Error())
		}
		return nil
	})
}

func (c *Client) DeleteAssistantWorkflow(ctx context.Context, assistantId uint32, workflowId string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		cond := sqlopt.SQLOptions(
			sqlopt.WithID(assistantId),
			sqlopt.WithWorkflowID(workflowId),
		).Apply(tx)
		// 2. 检查数据是否存在
		var exists int64
		if err := cond.Model(&model.AssistantWorkflow{}).Count(&exists).Error; err != nil {
			return toErrStatus("assistant_workflow_delete", err.Error())
		}
		if exists == 0 {
			return toErrStatus("assistant_workflow_delete", "workflow not exist")
		}

		// 删除Workflow
		if err := cond.Delete(&model.AssistantWorkflow{}).Error; err != nil {
			return toErrStatus("assistant_workflow_delete", err.Error())
		}

		// 查询是否还有其他Workflow
		var count int64
		if err := sqlopt.WithWorkflowID(workflowId).Apply(tx).Model(&model.AssistantWorkflow{}).Count(&count).Error; err != nil {
			return toErrStatus("assistant_workflows_count", err.Error())
		}

		// 如果没有其他Workflow，更新Assistant的HasWorkflow字段
		if count == 0 {
			if err := sqlopt.WithID(assistantId).Apply(tx).Model(&model.Assistant{}).Update("has_workflow", false).Error; err != nil {
				return toErrStatus("assistant_update", err.Error())
			}
		}

		return nil
	})
}
