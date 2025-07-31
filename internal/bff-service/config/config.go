package config

import (
	"github.com/UnicomAI/wanwu/pkg/i18n"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Config struct {
	Server            ServerConfig            `json:"server" mapstructure:"server"`
	Log               LogConfig               `json:"log" mapstructure:"log"`
	JWT               JWTConfig               `json:"jwt" mapstructure:"jwt"`
	Decrypt           DecryptPasswd           `json:"decrypt-passwd" mapstructure:"decrypt-passwd"`
	I18n              i18n.Config             `json:"i18n" mapstructure:"i18n"`
	AssistantTemplate AssistantTemplateConfig `json:"assistant-template" mapstructure:"assistant-template"`
	CustomInfo        CustomInfoConfig        `json:"custom-info" mapstructure:"custom-info"`
	DocCenter         DocCenterConfig         `json:"doc-center" mapstructure:"doc-center"`
	// middleware
	Minio minio.Config `json:"minio" mapstructure:"minio"`
	// microservice
	Iam       ServiceConfig         `json:"iam" mapstructure:"iam"`
	Model     ModelConfig           `json:"model" mapstructure:"model"`
	MCP       ServiceConfig         `json:"mcp" mapstructure:"mcp"`
	App       ServiceConfig         `json:"app" mapstructure:"app"`
	Knowledge ServiceConfig         `json:"knowledge" mapstructure:"knowledge"`
	Rag       ServiceConfig         `json:"rag" mapstructure:"rag"`
	Assistant ServiceConfig         `json:"assistant" mapstructure:"assistant"`
	WorkFlow  WorkFlowServiceConfig `json:"workflow" mapstructure:"workflow"`
	Agent     AgentServiceConfig    `json:"agent" mapstructure:"agent"`
}

type ServerConfig struct {
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	ExternalIP   string `json:"external_ip" mapstructure:"external_ip"`
	ExternalPort int    `json:"external_port" mapstructure:"external_port"`
	WebBaseUrl   string `json:"web_base_url" mapstructure:"web_base_url"`
	ApiBaseUrl   string `json:"api_base_url" mapstructure:"api_base_url"`
	CallbackUrl  string `json:"callback_url" mapstructure:"callback_url"`
}

type ModelConfig struct {
	Host        string `json:"host" mapstructure:"host"`
	OcrFilePath string `json:"ocr_test_file_path" mapstructure:"ocr_test_file_path"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

type JWTConfig struct {
	SigningKey string `json:"signing-key" mapstructure:"signing-key"`
}

type DecryptPasswd struct {
	IV  string `json:"iv" mapstructure:"iv"`
	Key string `json:"key" mapstructure:"key"`
}

type ServiceConfig struct {
	Host string `json:"host" mapstructure:"host"`
}

type WorkFlowServiceConfig struct {
	Endpoint                string `json:"host" mapstructure:"host"`
	WorkFlowListUri         string `json:"workflow_list_uri" mapstructure:"workflow_list_uri"`
	WorkFlowListUriInternal string `json:"workflow_list_uri_internal" mapstructure:"workflow_list_uri_internal"`
	DeleteWorkFlowUri       string `json:"delete_workflow_uri" mapstructure:"delete_workflow_uri"`
	PublishWorkFlowUri      string `json:"publish_workflow_uri" mapstructure:"publish_workflow_uri"`
	UnPublishWorkFlowUri    string `json:"unpublish_workflow_uri" mapstructure:"unpublish_workflow_uri"`
}

type AgentServiceConfig struct {
	Endpoint       string    `json:"host" mapstructure:"host"`
	UploadMinioUri UriConfig `json:"upload_minio" mapstructure:"upload_minio"`
}

type UriConfig struct {
	Port string `json:"port" mapstructure:"port"`
	Uri  string `json:"uri" mapstructure:"uri"`
}

type AssistantTemplateConfig struct {
	ConfigPath string `json:"configPath" mapstructure:"configPath"`
}

type DocCenterConfig struct {
	DocPath string `json:"doc_path" mapstructure:"doc_path"`
}

type CustomInfoConfig struct {
	Login CustomLogin `json:"login" mapstructure:"login"`
	Home  CustomHome  `json:"home" mapstructure:"home"`
	Tab   CustomTab   `json:"tab" mapstructure:"tab"`
	About CustomAbout `json:"about" mapstructure:"about"`
}

type CustomLogin struct {
	BackgroundPath   string `json:"background_path" mapstructure:"background_path"`
	LoginButtonColor string `json:"login_button_color" mapstructure:"login_button_color"`
	WelcomeText      string `json:"welcome_text" mapstructure:"welcome_text"`
	PlatformDesc     string `json:"platform_desc" mapstructure:"platform_desc"`
}

type CustomHome struct {
	LogoPath string `json:"logo_path" mapstructure:"logo_path"`
	Title    string `json:"title" mapstructure:"title"`
}

type CustomTab struct {
	TabTitle    string `json:"tab_title" mapstructure:"tab_title"`
	TabLogoPath string `json:"tab_logo_path" mapstructure:"tab_logo_path"`
}

type CustomAbout struct {
	LogoPath  string `json:"logo_path" mapstructure:"logo_path"`
	Version   string `json:"version" mapstructure:"version"`
	Copyright string `json:"copyright" mapstructure:"copyright"`
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
