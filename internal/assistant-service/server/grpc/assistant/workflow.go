package assistant

import (
	"context"
	"strconv"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AssistantWorkFlowCreate 添加workFlow
func (s *Service) AssistantWorkFlowCreate(ctx context.Context, req *assistant_service.AssistantWorkFlowCreateReq) (*emptypb.Empty, error) {
	// 组装model参数
	assistantID, err := strconv.ParseUint(req.AssistantId, 10, 32)
	if err != nil {
		return nil, err
	}

	workflow := &model.AssistantWorkflow{
		AssistantId: uint32(assistantID),
		APISchema:   req.Schema,
		UserId:      req.Identity.UserId,
		OrgId:       req.Identity.OrgId,
	}

	// 调用client方法创建WorkFlow（在事务中创建WorkFlow并更新Assistant）
	if status := s.cli.CreateAssistantWorkflow(ctx, workflow); status != nil {
		return nil, errStatus(errs.Code_AssistantWorkflowErr, status)
	}

	return &emptypb.Empty{}, nil
}

// AssistantWorkFlowDelete 删除workFlow
func (s *Service) AssistantWorkFlowDelete(ctx context.Context, req *assistant_service.AssistantWorkFlowDeleteReq) (*emptypb.Empty, error) {
	// 转换ID
	workflowID, err := strconv.ParseUint(req.WorkFlowId, 10, 32)
	if err != nil {
		return nil, err
	}

	// 调用client方法删除WorkFlow（在事务中删除WorkFlow并更新Assistant）
	if status := s.cli.DeleteAssistantWorkflow(ctx, uint32(workflowID)); status != nil {
		return nil, errStatus(errs.Code_AssistantWorkflowErr, status)
	}

	return &emptypb.Empty{}, nil
}

// AssistantWorkFlowEnableSwitch WorkFlow开关
func (s *Service) AssistantWorkFlowEnableSwitch(ctx context.Context, req *assistant_service.AssistantWorkFlowEnableSwitchReq) (*emptypb.Empty, error) {
	// 转换ID
	workflowID, err := strconv.ParseUint(req.WorkFlowId, 10, 32)
	if err != nil {
		return nil, err
	}

	// 先获取现有WorkFlow信息
	existingWorkflow, status := s.cli.GetAssistantWorkflow(ctx, uint32(workflowID))
	if status != nil {
		return nil, errStatus(errs.Code_AssistantWorkflowErr, status)
	}

	existingWorkflow.Enable = !existingWorkflow.Enable
	if status := s.cli.UpdateAssistantWorkflow(ctx, existingWorkflow); status != nil {
		return nil, errStatus(errs.Code_AssistantWorkflowErr, status)
	}

	return &emptypb.Empty{}, nil
}
