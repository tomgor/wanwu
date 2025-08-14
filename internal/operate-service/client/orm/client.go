package orm

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/operate-service/client/model"
	"gorm.io/gorm"
)

type SystemCustomKey string
type SystemCustomMode string

const (
	SystemCustomTabKey   SystemCustomKey = "system_custom_tab"
	SystemCustomLoginKey SystemCustomKey = "system_custom_login"
	SystemCustomHomeKey  SystemCustomKey = "system_custom_home"
)
const (
	SystemCustomModeLight SystemCustomMode = "light"
	SystemCustomModeDark  SystemCustomMode = "dark"
)

type Client struct {
	db *gorm.DB
}

func NewClient(db *gorm.DB) (*Client, error) {
	// auto migrate
	if err := db.AutoMigrate(
		model.SystemCustom{},
	); err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

func toErrStatus(key string, args ...string) *err_code.Status {
	return &err_code.Status{
		TextKey: key,
		Args:    args,
	}
}

type SystemCustom struct {
	Login LoginConfig `json:"login"` // 登录页配置
	Tab   TabConfig   `json:"tab"`   // 标签页配置
	Home  HomeConfig  `json:"home"`  // 首页配置
}

type LoginConfig struct {
	LoginBgPath string `json:"loginBgPath"` // 登录页背景图路径
	WelcomeText string `json:"welcomeText"` // 登录页欢迎词
	ButtonColor string `json:"buttonColor"` // 登录按钮颜色
}

type TabConfig struct {
	LogoPath string `json:"logoPath"` // 标签页logo路径
	Title    string `json:"title"`    // 标签页标题
}

type HomeConfig struct {
	LogoPath string `json:"logoPath"` // 平台logo路径
	Name     string `json:"name"`     // 平台名称
	BgColor  string `json:"bgColor"`  // 平台背景颜色
}
