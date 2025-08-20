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
	workflow, err := parseWorkFlowApiInfo(req)
	if err != nil {
		return nil, err
	}

	// 调用client方法创建WorkFlow（在事务中创建WorkFlow并更新Assistant）
	if status := s.cli.CreateAssistantWorkflow(ctx, workflow); status != nil {
		return nil, errStatus(errs.Code_AssistantWorkflowErr, status)
	}

	return &emptypb.Empty{}, nil
}

func parseWorkFlowApiInfo(req *assistant_service.AssistantWorkFlowCreateReq) (*model.AssistantWorkflow, error) {
	userId, orgId := req.Identity.UserId, req.Identity.OrgId
	assistantID, err := strconv.ParseUint(req.AssistantId, 10, 32)
	if err != nil {
		return nil, err
	}
	workFlowId := req.WorkFlowId
	workFlow := &model.AssistantWorkflow{
		WorkflowId:  workFlowId,
		AssistantId: uint32(assistantID),
		Enable:      true,
		UserId:      userId,
		OrgId:       orgId,
	}

	return workFlow, nil
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
