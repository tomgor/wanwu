package minio

import (
	"context"
	"fmt"
)

var (
	_minioKnowledge *client
)

func InitKnowledge(ctx context.Context, cfg Config, initBucketName string) error {
	if _minioKnowledge != nil {
		return fmt.Errorf("minio knowledge client already init")
	}
	c, err := newClient(cfg)
	if err != nil {
		return err
	}
	_minioKnowledge = c
	if err = _minioKnowledge.CreateBucketIfNotExist(ctx, initBucketName); err != nil {
		return err
	}
	return err
}

func Knowledge() *client {
	return _minioKnowledge
}
