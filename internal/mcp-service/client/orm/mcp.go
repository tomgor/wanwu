package orm

import (
	"context"
	"database/sql"
	"errors"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_client "github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) GetMCP(ctx context.Context, tab *mcp_client.MCPModel) (*mcp_client.MCPModel, *errs.Status) {
	info := &mcp_client.MCPModel{}
	err := sqlopt.SQLOptions(
		sqlopt.WithMcpSquareId(tab.McpSquareId),
		sqlopt.WithID(tab.ID),
		sqlopt.WithOrgID(tab.OrgID),
		sqlopt.WithUserID(tab.UserID),
	).Apply(c.db).WithContext(ctx).First(info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, toErrStatus("mcp_get_err", err.Error())
	}
	return info, nil
}

func (c *Client) CreateMCP(ctx context.Context, tab *mcp_client.MCPModel) *errs.Status {
	var err error
	db := c.db.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable}).WithContext(ctx)
	defer func() {
		if err == nil {
			db.Commit()
			return
		}
		db.Rollback()
	}()
	// 先查询是否已存在相同的记录
	if err = sqlopt.SQLOptions(
		sqlopt.WithName(tab.Name),
		sqlopt.WithOrgID(tab.OrgID),
		sqlopt.WithUserID(tab.UserID),
	).Apply(db).Select("id").First(&mcp_client.MCPModel{}).Error; err == nil {
		return toErrStatus("mcp_create_err", "mcp server with same name exist")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return toErrStatus("mcp_create_err", err.Error())
	}

	if err = db.Create(tab).Error; err != nil {
		return toErrStatus("mcp_create_err", err.Error())
	}
	return nil
}

func (c *Client) DeleteMCP(ctx context.Context, tab *mcp_client.MCPModel) *errs.Status {
	var existing mcp_client.MCPModel
	if err := sqlopt.SQLOptions(
		sqlopt.WithID(tab.ID),
	).Apply(c.db).WithContext(ctx).First(&existing).Error; err != nil {
		return toErrStatus("mcp_delete_err", err.Error())
	}
	if err := c.db.WithContext(ctx).Delete(existing).Error; err != nil {
		return toErrStatus("mcp_delete_err", err.Error())
	}
	return nil
}

func (c *Client) ListMCPs(ctx context.Context, tab *mcp_client.MCPModel) ([]*mcp_client.MCPModel, int64, *errs.Status) {
	var count int64
	var mcpInfos []*mcp_client.MCPModel
	db := sqlopt.SQLOptions(
		sqlopt.WithOrgID(tab.OrgID),
		sqlopt.WithUserID(tab.UserID),
		sqlopt.LikeName(tab.Name),
	).Apply(c.db.WithContext(ctx))

	if err := db.Order("created_at DESC").Find(&mcpInfos).Error; err != nil {
		return nil, 0, toErrStatus("mcp_list_mcps_err", err.Error())
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, toErrStatus("mcp_list_mcps_err", err.Error())
	}
	return mcpInfos, count, nil
}
