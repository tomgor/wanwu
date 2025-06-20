package mcp

import (
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client"
)

type Service struct {
	mcp_service.UnimplementedMCPServiceServer
	cli client.IClient
}

func NewService(cli client.IClient) *Service {
	return &Service{
		cli: cli,
	}
}
