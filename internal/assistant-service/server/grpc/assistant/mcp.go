package assistant

import (
	"context"
	"strconv"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AssistantMCPCreate 添加mcp
func (s *Service) AssistantMCPCreate(ctx context.Context, req *assistant_service.AssistantMCPCreateReq) (*emptypb.Empty, error) {
	assistantId, err := strconv.ParseUint(req.AssistantId, 10, 32)
	if err != nil {
		return nil, err
	}

	if status := s.cli.CreateAssistantMCP(ctx, &model.AssistantMCP{
		AssistantId: uint32(assistantId),
		MCPId:       req.McpId,
	}); status != nil {
		return nil, errStatus(errs.Code_AssistantMCPErr, status)
	}

	return &emptypb.Empty{}, nil
}

// AssistantMCPDelete 删除mcp
func (s *Service) AssistantMCPDelete(ctx context.Context, req *assistant_service.AssistantMCPDeleteReq) (*emptypb.Empty, error) {
	// 这里的 MCPID 实为 ID
	id, err := strconv.ParseUint(req.McpId, 10, 32)
	if err != nil {
		return nil, err
	}

	if status := s.cli.DeleteAssistantMCP(ctx, uint32(id)); status != nil {
		return nil, errStatus(errs.Code_AssistantMCPErr, status)
	}
	return &emptypb.Empty{}, nil
}

// AssistantMCPEnableSwitch mcp开关
func (s *Service) AssistantMCPEnableSwitch(ctx context.Context, req *assistant_service.AssistantMCPEnableSwitchReq) (*emptypb.Empty, error) {
	// 这里的 MCPID 实为 ID
	id, err := strconv.ParseUint(req.McpId, 10, 32)
	if err != nil {
		return nil, err
	}

	// 先获取现有MCP信息
	query := map[string]interface{}{
		"id": id,
	}

	existingMCP, status := s.cli.GetAssistantMCP(ctx, query)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantMCPErr, status)
	}

	existingMCP.Enable = !existingMCP.Enable
	if status := s.cli.UpdateAssistantMCP(ctx, existingMCP); status != nil {
		return nil, errStatus(errs.Code_AssistantMCPErr, status)
	}

	return &emptypb.Empty{}, nil
}
func (s *Service) AssistantMCPGetList(ctx context.Context, req *assistant_service.AssistantMCPGetListReq) (*assistant_service.AssistantMCPList, error) {
	assistantID, err := strconv.ParseUint(req.AssistantId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := map[string]interface{}{
		"assistant_id": assistantID,
	}

	mcpList, status := s.cli.GetAssistantMCPList(ctx, query)
	if status != nil {
		return nil, errStatus(errs.Code_AssistantMCPErr, status)
	}

	assistantMCPInfos := make([]*assistant_service.AssistantMCPInfo, len(mcpList))
	for i, mcp := range mcpList {
		assistantMCPInfos[i] = &assistant_service.AssistantMCPInfo{
			Id:             mcp.ID,
			AssistantMcpId: mcp.AssistantId,
			McpId:          mcp.MCPId,
			Enable:         mcp.Enable,
			UserId:         mcp.UserId,
			OrgId:          mcp.OrgId,
			CreatedAt:      mcp.CreatedAt,
			UpdatedAt:      mcp.UpdatedAt,
		}
	}

	return &assistant_service.AssistantMCPList{
		AssistantMCPInfos: assistantMCPInfos,
	}, nil

}
