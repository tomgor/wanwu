package mcp

import (
	"context"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/internal/mcp-service/config"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Service) GetSquareMCP(ctx context.Context, req *mcp_service.GetSquareMCPReq) (*mcp_service.SquareMCPDetail, error) {
	if req.McpSquareId == "" {
		return nil, nil
	}

	for _, server := range config.Cfg().Mcp {
		if server.McpSquareId != req.McpSquareId {
			continue
		}

		hasCustom, err := s.checkCustomMCP(ctx, req)
		if err != nil {
			return nil, errStatus(errs.Code_MCPGetSquareMCPErr, err)
		}

		return &mcp_service.SquareMCPDetail{
			Info:  s.buildSquareMCPInfo(&server),
			Intro: s.buildSquareMCPIntro(&server),
			Tool: &mcp_service.MCPTools{
				SseUrl:    server.SseUrl,
				HasCustom: hasCustom,
				Tools:     s.convertTools(server.Tools),
			},
		}, nil
	}
	return nil, nil
}

func (s *Service) checkCustomMCP(ctx context.Context, req *mcp_service.GetSquareMCPReq) (bool, *errs.Status) {
	mcpInfo, err := s.cli.GetMCP(ctx, &model.MCPModel{
		McpSquareId: req.McpSquareId,
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})
	if err != nil {
		return false, toErrStatus("mcp_check_custom_err")
	}
	return mcpInfo != nil, nil
}

func (s *Service) buildSquareMCPInfo(server *config.McpConfig) *mcp_service.SquareMCPInfo {
	return &mcp_service.SquareMCPInfo{
		Name:        server.Name,
		Desc:        server.Desc,
		Category:    server.Category,
		McpSquareId: server.McpSquareId,
		From:        server.From,
		AvatarPath:  server.Avatar,
	}
}

func (s *Service) buildSquareMCPIntro(server *config.McpConfig) *mcp_service.SquareMCPIntro {
	mdContent, err := os.ReadFile(filepath.Join(ConfigDir, server.Detail))
	if err != nil {
		log.Errorf("read square mcp detail failed, path: %s, err: %v", server.Detail, err)
	}
	return &mcp_service.SquareMCPIntro{
		Summary:  server.Summary,
		Detail:   string(mdContent),
		Feature:  server.Feature,
		Manual:   server.Manual,
		Scenario: server.Scenario,
	}
}

func (s *Service) convertTools(tools []config.McpToolConfig) []*mcp_service.MCPTool {
	result := make([]*mcp_service.MCPTool, 0, len(tools))
	for _, tool := range tools {
		result = append(result, &mcp_service.MCPTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: s.convertInputSchema(&tool.InputSchema),
		})
	}
	return result
}

func (s *Service) convertInputSchema(schema *config.McpInputSchemaConfig) *mcp_service.MCPToolInputSchema {
	if schema == nil {
		return nil
	}

	properties := make(map[string]*mcp_service.MCPToolInputSchemaValue)
	for _, prop := range schema.Properties {
		properties[prop.Field] = &mcp_service.MCPToolInputSchemaValue{
			Type:        prop.Type,
			Description: prop.Description,
		}
	}

	return &mcp_service.MCPToolInputSchema{
		Type:       schema.Type,
		Required:   schema.Required,
		Properties: properties,
	}
}

func (s *Service) GetSquareMCPList(ctx context.Context, req *mcp_service.GetSquareMCPListReq) (*mcp_service.SquareMCPList, error) {
	var resMcpSquareServers []*mcp_service.SquareMCPInfo

	for _, server := range config.Cfg().Mcp {
		shouldAppend := true

		if req.Name != "" && !strings.Contains(server.Name, req.Name) {
			shouldAppend = false
		}

		if req.Category != "" && !strings.Contains(server.Category, req.Category) {
			shouldAppend = false
		}

		if shouldAppend {
			resMcpSquareServers = append(resMcpSquareServers, &mcp_service.SquareMCPInfo{
				McpSquareId: server.McpSquareId,
				AvatarPath:  server.Avatar,
				Name:        server.Name,
				Desc:        server.Desc,
				Category:    server.Category,
				From:        server.From,
			})
		}
	}

	return &mcp_service.SquareMCPList{Infos: resMcpSquareServers}, nil
}

func (s *Service) CreateCustomMCP(ctx context.Context, req *mcp_service.CreateCustomMCPReq) (*emptypb.Empty, error) {
	if err := s.cli.CreateMCP(ctx, &model.MCPModel{
		Name:        req.Name,
		McpSquareId: req.McpSquareId,
		Desc:        req.Desc,
		From:        req.From,
		SseUrl:      req.SseUrl,
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	}); err != nil {
		return nil, errStatus(errs.Code_MCPCreateCustomMCPErr, err)
	}
	return nil, nil
}

func (s *Service) GetCustomMCP(ctx context.Context, req *mcp_service.GetCustomMCPReq) (*mcp_service.CustomMCPDetail, error) {
	mcpInfo, err := s.cli.GetMCP(ctx, &model.MCPModel{
		ID: util.MustU32(req.McpId),
	})
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomMCPErr, err)
	}

	res := &mcp_service.CustomMCPDetail{
		McpId:  strconv.Itoa(int(mcpInfo.ID)),
		SseUrl: mcpInfo.SseUrl,
	}

	if mcpInfo.McpSquareId == "" {
		res.Info = &mcp_service.SquareMCPInfo{
			Name: mcpInfo.Name,
			Desc: mcpInfo.Desc,
			From: mcpInfo.From,
		}
		return res, nil
	}

	for _, server := range config.Cfg().Mcp {
		if server.McpSquareId != mcpInfo.McpSquareId {
			continue
		}
		res.Info = s.buildSquareMCPInfo(&server)
		res.Intro = s.buildSquareMCPIntro(&server)
	}

	return res, nil
}

func (s *Service) DeleteCustomMCP(ctx context.Context, req *mcp_service.DeleteCustomMCPReq) (*emptypb.Empty, error) {
	if err := s.cli.DeleteMCP(ctx, &model.MCPModel{
		ID: util.MustU32(req.McpId),
	}); err != nil {
		return nil, errStatus(errs.Code_MCPDeleteCustomMCPErr, err)
	}

	return nil, nil
}

func (s *Service) GetCustomMCPList(ctx context.Context, req *mcp_service.GetCustomMCPListReq) (*mcp_service.CustomMCPList, error) {
	mcpInfoList, _, err := s.cli.ListMCPs(ctx, &model.MCPModel{
		Name: req.Name,
		PublicModel: model.PublicModel{
			OrgID:  req.OrgId,
			UserID: req.UserId,
		},
	})
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomMCPListErr, err)
	}

	infos := make([]*mcp_service.CustomMCPInfo, 0, len(mcpInfoList))
	for _, mcpInfo := range mcpInfoList {
		info := &mcp_service.CustomMCPInfo{
			McpId:  strconv.Itoa(int(mcpInfo.ID)),
			SseUrl: mcpInfo.SseUrl,
			Info:   s.buildMCPInfo(mcpInfo),
		}
		infos = append(infos, info)
	}

	return &mcp_service.CustomMCPList{Infos: infos}, nil
}

func (s *Service) buildMCPInfo(mcpInfo *model.MCPModel) *mcp_service.SquareMCPInfo {
	if mcpInfo.McpSquareId == "" {
		return &mcp_service.SquareMCPInfo{
			Name: mcpInfo.Name,
			Desc: mcpInfo.Desc,
			From: mcpInfo.From,
		}
	}

	for _, server := range config.Cfg().Mcp {
		if server.McpSquareId == mcpInfo.McpSquareId {
			return s.buildSquareMCPInfo(&server)
		}
	}
	return nil
}

func (s *Service) GetMCPAvatar(ctx context.Context, req *mcp_service.GetMCPAvatarReq) (*mcp_service.MCPAvatar, error) {
	if req.AvatarPath == "" {
		return nil, errStatus(errs.Code_MCPGetMCPAvatarErr, toErrStatus("mcp_get_mcp_avatar_err", "avatar path is empty"))
	}
	filePath := filepath.Join(ConfigDir, req.AvatarPath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetMCPAvatarErr, toErrStatus("mcp_get_mcp_avatar_err", err.Error()))
	}

	_, fileName := filepath.Split(filePath)

	return &mcp_service.MCPAvatar{
		FileName: fileName,
		Data:     data,
	}, nil
}
