package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type ToolConfig struct {
	ToolSquareId     string          `json:"tool_square_id" mapstructure:"tool_square_id"`
	Name             string          `json:"name" mapstructure:"name"`
	Desc             string          `json:"desc" mapstructure:"desc"`
	AvatarPath       string          `json:"avatar_path" mapstructure:"avatar_path"`
	Detail           string          `json:"detail" mapstructure:"detail"`
	Tags             string          `json:"tags" mapstructure:"tags"`
	Tools            []McpToolConfig `json:"tools" mapstructure:"tools"`
	Type             string          `json:"type" mapstructure:"type"`
	AuthType         string          `json:"auth_type" mapstructure:"auth_type"`
	CustomHeaderName string          `json:"custom_header_name" mapstructure:"custom_header_name"`
	ApiKey           string          `json:"api_key" mapstructure:"api_key"`
	Schema           string          `json:"schema" mapstructure:"schema"`
	NeedApiKeyInput  bool            `json:"need_api_key_input" mapstructure:"need_api_key_input"`
}

func (tool *ToolConfig) load() error {
	avatarPath := filepath.Join(ConfigDir, tool.AvatarPath)
	if _, err := os.ReadFile(avatarPath); err != nil {
		return fmt.Errorf("load tool %v avatar path %v err: %v", tool.ToolSquareId, avatarPath, err)
	}
	return nil
}
