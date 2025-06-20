package client

import (
	"context"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
)

type IClient interface {
	GetMCP(ctx context.Context, req *model.MCPModel) (*model.MCPModel, *errs.Status)
	CreateMCP(ctx context.Context, req *model.MCPModel) *errs.Status
	DeleteMCP(ctx context.Context, req *model.MCPModel) *errs.Status
	ListMCPs(ctx context.Context, req *model.MCPModel) ([]*model.MCPModel, int64, *errs.Status)
}
