package mcp

import (
	"context"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateCustomTool(ctx context.Context, req *mcp_service.CreateCustomToolReq) (*emptypb.Empty, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "identity is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "apiAuth is empty"))
	}
	customToolId := util.GenUUID()
	if err := s.cli.CreateCustomTool(ctx, &model.CustomTool{
		CustomToolId:     customToolId,
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
		UserID:           req.Identity.UserId,
		OrgID:            req.Identity.OrgId,
	}); err != nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *Service) GetCustomToolInfo(ctx context.Context, req *mcp_service.GetCustomToolInfoReq) (*mcp_service.GetCustomToolInfoResp, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, toErrStatus("mcp_get_custom_tool_info_err", "customToolId is empty"))
	}
	info, err := s.cli.GetCustomTool(ctx, req.CustomToolId)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, err)
	}
	return &mcp_service.GetCustomToolInfoResp{
		CustomToolId:  info.CustomToolId,
		Name:          info.Name,
		Description:   info.Description,
		Schema:        info.Schema,
		PrivacyPolicy: info.PrivacyPolicy,
		ApiAuth: &mcp_service.ApiAuthWebRequest{
			Type:             info.Type,
			ApiKey:           info.APIKey,
			AuthType:         info.AuthType,
			CustomHeaderName: info.CustomHeaderName,
		},
	}, nil
}
func (s *Service) GetCustomToolList(ctx context.Context, req *mcp_service.GetCustomToolListReq) (*mcp_service.GetCustomToolListResp, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "identity is empty"))
	}
	infos, err := s.cli.ListCustomTools(ctx, req.Identity.OrgId, req.Identity.UserId, req.Name)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: info.CustomToolId,
			Name:         info.Name,
			Description:  info.Description,
		})
	}
	return &mcp_service.GetCustomToolListResp{
		List: list,
	}, nil
}
func (s *Service) UpdateCustomTool(ctx context.Context, req *mcp_service.UpdateCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "customToolId is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "apiAuth is empty"))
	}
	if err := s.cli.UpdateCustomTool(ctx, &model.CustomTool{
		CustomToolId:     req.CustomToolId,
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
	}); err != nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}
func (s *Service) DeleteCustomTool(ctx context.Context, req *mcp_service.DeleteCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, toErrStatus("mcp_delete_custom_tool_err", "customToolId is empty"))
	}
	if err := s.cli.DeleteCustomTool(ctx, req.CustomToolId); err != nil {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}
