package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"

	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(ctx context.Context, db *gorm.DB) (*Client, error) {
	// auto migrate
	if err := db.AutoMigrate(
		model.MCPModel{},
	); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}
