package mcp

import (
	"context"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetSquareMCP(ctx context.Context, req *mcp_service.GetSquareMCPReq) (*mcp_service.SquareMCPDetail, error) {
	//McpSquareServers := config.Cfg().Mcp
	//if req.McpSquareId == "" {
	//	return nil, nil
	//}
	//for _, mcpSquareServer := range McpSquareServers {
	//	if mcpSquareServer.McpSquareId == req.McpSquareId {
	//		log.Infof("tools: %v", mcpSquareServer.Tools[0].InputSchema.Properties[0])
	//		resMcpTools := make([]*mcp_service.MCPTool, 0)
	//		var resMcpTool *mcp_service.MCPTool
	//		for _, tool := range mcpSquareServer.Tools {
	//			resMcpTool = &mcp_service.MCPTool{}
	//			resMcpTool.InputSchema = &mcp_service.MCPToolInputSchema{}
	//
	//			resMcpTool.Name = tool.Name
	//			resMcpTool.Description = tool.Description
	//			resMcpTool.InputSchema.Type = tool.InputSchema.Type
	//			resMcpTool.InputSchema.Required = tool.InputSchema.Required
	//			resMcpTool.InputSchema.Properties = make(map[string]*mcp_service.MCPToolInputSchemaValue)
	//			for _, propertie := range tool.InputSchema.Properties {
	//				resMcpTool.InputSchema.Properties[propertie.Field] = &mcp_service.MCPToolInputSchemaValue{
	//					Type:        propertie.Type,
	//					Description: propertie.Description,
	//				}
	//			}
	//			resMcpTools = append(resMcpTools, resMcpTool)
	//		}
	//		return &mcp_service.SquareMCPDetail{
	//			Info: &mcp_service.SquareMCPInfo{
	//				Name:        mcpSquareServer.Name,
	//				Desc:        mcpSquareServer.Desc,
	//				Category:    mcpSquareServer.Category,
	//				McpSquareId: mcpSquareServer.McpSquareId,
	//				From:        mcpSquareServer.From,
	//				AvatarPath:  mcpSquareServer.Avatar,
	//			},
	//			Intro: &mcp_service.SquareMCPIntro{
	//				Summary:  mcpSquareServer.Summary,
	//				Detail:   mcpSquareServer.Detail,
	//				Feature:  mcpSquareServer.Feature,
	//				Manual:   mcpSquareServer.Manual,
	//				Scenario: mcpSquareServer.Scenario,
	//			},
	//			//查询库里是否有该McpSquareId, 有的话HasCustom=true，否则为false
	//			Tool: &mcp_service.MCPTools{
	//				SseUrl:    mcpSquareServer.SseUrl,
	//				HasCustom: false,
	//				Tools:     resMcpTools,
	//			},
	//		}, nil
	//	}
	//}
	return nil, nil
}

func (s *Service) GetSquareMCPList(ctx context.Context, req *mcp_service.GetSquareMCPListReq) (*mcp_service.SquareMCPList, error) {
	return nil, nil
}

func (s *Service) CreateCustomMCP(ctx context.Context, req *mcp_service.CreateCustomMCPReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) GetCustomMCP(ctx context.Context, req *mcp_service.GetCustomMCPReq) (*mcp_service.CustomMCPDetail, error) {
	return nil, nil
}

func (s *Service) DeleteCustomMCP(ctx context.Context, req *mcp_service.DeleteCustomMCPReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) GetCustomMCPList(ctx context.Context, req *mcp_service.GetCustomMCPListReq) (*mcp_service.CustomMCPList, error) {
	return nil, nil
}

func (s *Service) GetMCPAvatar(ctx context.Context, req *mcp_service.GetMCPAvatarReq) (*mcp_service.MCPAvatar, error) {
	return nil, nil
}
