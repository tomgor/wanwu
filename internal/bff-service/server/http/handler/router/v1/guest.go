package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerGuest(apiV1 *gin.RouterGroup) {
	apiV1.Static("/static", "./configs/microservice/bff-service/static")
	apiV1.Static("/cache", "./cache")

	mid.Sub("guest").Reg(apiV1, "/base/register/email", http.MethodPost, v1.RegisterByEmail, "邮箱注册用户")
	mid.Sub("guest").Reg(apiV1, "/base/register/email/code", http.MethodPost, v1.ResgisterSendEmailCode, "邮箱注册验证码发送")

	mid.Sub("guest").Reg(apiV1, "/base/password/email/code", http.MethodPost, v1.ResetPasswordSendEmailCode, "重置密码邮箱验证码发送")
	mid.Sub("guest").Reg(apiV1, "/base/password/email", http.MethodPost, v1.ResetPasswordByEmail, "邮箱重置密码")

	mid.Sub("guest").Reg(apiV1, "/base/login", http.MethodPost, v1.Login, "用户登录")
	mid.Sub("guest").Reg(apiV1, "/base/captcha", http.MethodGet, v1.GetCaptcha, "获取验证码")
	mid.Sub("guest").Reg(apiV1, "/base/custom", http.MethodGet, v1.GetLogoCustomInfo, "自定义logo和title")
	mid.Sub("guest").Reg(apiV1, "/base/language/select", http.MethodGet, v1.GetLanguageSelect, "获取语言列表（用于下拉选择）")

	mid.Sub("guest").Reg(apiV1, "/base/simple-sso", http.MethodPost, v1.SimpleSSO, "简单用户单点登陆")
}
