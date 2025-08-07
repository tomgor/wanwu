package orm

import (
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(db *gorm.DB) (*Client, error) {
	// auto migrate
	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}
