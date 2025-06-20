package config

import (
	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Mcp struct {
}
type Config struct {
	Server ServerConfig `json:"server" mapstructure:"server"`
	Log    LogConfig    `json:"log" mapstructure:"log"`
	DB     db.Config    `json:"db" mapstructure:"db"`
	Mcp    []McpConfig  `json:"mcp" mapstructure:"mcp"`
}
type McpConfig struct {
	McpSquareId string          `json:"mcp_square_id" mapstructure:"mcp_square_id"`
	Name        string          `json:"name" mapstructure:"name"`
	Category    string          `json:"category" mapstructure:"category"`
	Desc        string          `json:"desc" mapstructure:"desc"`
	From        string          `json:"from" mapstructure:"from"`
	Avatar      string          `json:"avatar" mapstructure:"avatar"`
	Detail      string          `json:"detail" mapstructure:"detail"`
	Feature     string          `json:"feature" mapstructure:"feature"`
	Manual      string          `json:"manual" mapstructure:"manual"`
	Scenario    string          `json:"scenario" mapstructure:"scenario"`
	SseUrl      string          `json:"sse_url" mapstructure:"sse_url"`
	Summary     string          `json:"summary" mapstructure:"summary"`
	Tools       []McpToolConfig `json:"tools" mapstructure:"tools"`
}

type McpToolConfig struct {
	Name        string               `json:"name" mapstructure:"name"`
	Description string               `json:"description" mapstructure:"description"`
	InputSchema McpInputSchemaConfig `json:"input_schema" mapstructure:"input_schema"`
}

type McpInputSchemaConfig struct {
	Type       string                `json:"type" mapstructure:"type"`
	Required   []string              `json:"required" mapstructure:"required"`
	Properties []McpPropertiesConfig `json:"properties" mapstructure:"properties"`
}

type McpPropertiesConfig struct {
	Field       string `json:"field" mapstructure:"field"`
	Type        string `json:"type" mapstructure:"type"`
	Description string `json:"description" mapstructure:"description"`
}

type ServerConfig struct {
	GrpcEndpoint   string `json:"grpc_endpoint" mapstructure:"grpc_endpoint"`
	MaxRecvMsgSize int    `json:"max_recv_msg_size" mapstructure:"max_recv_msg_size"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

type DBConfig struct {
	Name string `json:"name" mapstructure:"name"`
}

func LoadConfig(in string) error {
	_c = &Config{}
	return util.LoadConfig(in, _c)
}

func Cfg() *Config {
	if _c == nil {
		log.Panicf("cfg nil")
	}
	return _c
}
