package minio

import (
	"context"
	"fmt"
)

var (
	_minioAssistant *client
)

func InitAssistant(ctx context.Context, cfg Config, initBucketName string) error {
	if _minioAssistant != nil {
		return fmt.Errorf("minio assistat client already init")
	}
	c, err := newClient(cfg)
	if err != nil {
		return err
	}
	_minioAssistant = c
	if err = _minioAssistant.CreateBucketIfNotExist(ctx, initBucketName); err != nil {
		return err
	}
	return err
}
