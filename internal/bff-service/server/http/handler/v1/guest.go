package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

//	@title		AI Agent Productivity Platform API
//	@version	v0.0.1
//	@description.markdown
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization

//	@BasePath	/v1

// Login
//
//	@Tags		guest
//	@Summary	用户登录
//	@Accept		json
//	@Produce	json
//	@Param		X-Language	header		string			false	"语言"
//	@Param		data		body		request.Login	true	"用户名+密码"
//	@Success	200			{object}	response.Response{data=response.Login}
//	@Router		/base/login [post]
func Login(ctx *gin.Context) {
	var req request.Login
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.Login(ctx, &req, getLanguage(ctx))
	gin_util.Response(ctx, resp, err)
}

// GetCaptcha
//
//	@Tags		guest
//	@Summary	获取验证码
//	@Accept		json
//	@Produce	json
//	@Param		X-Language	header		string	false	"语言"
//	@Success	200			{object}	response.Response{data=response.Captcha}
//	@Router		/base/captcha [get]
func GetCaptcha(ctx *gin.Context) {
	resp, err := service.GetCaptcha(ctx,
		util.MD5([]byte(ctx.ClientIP()+ctx.GetHeader("User-Agent")+ctx.GetHeader("Date"))))
	gin_util.Response(ctx, resp, err)
}

// GetLogoCustomInfo
//
//	@Tags		guest
//	@Summary	自定义logo和title
//	@Produce	application/json
//	@Param		X-Language	header		string	false	"语言"
//	@Success	200			{object}	response.Response{data=response.LogoCustomInfo}
//	@Router		/base/custom [get]
func GetLogoCustomInfo(ctx *gin.Context) {
	resp, err := service.GetLogoCustomInfo(ctx, config.Cfg().CustomInfo.DefaultMode)
	gin_util.Response(ctx, resp, err)
}

// GetLanguageSelect
//
//	@Tags		guest
//	@Summary	获取语言列表（用于下拉选择）
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.LanguageSelect}
//	@Router		/base/language/select [get]
func GetLanguageSelect(ctx *gin.Context) {
	resp := service.GetLanguageSelect()
	gin_util.Response(ctx, resp, nil)
}

// RegisterByEmail
//
//	@Tags		guest
//	@Summary	用户邮箱注册
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.RegisterByEmail	true	"邮箱注册信息"
//	@Success	200		{object}	response.Response
//	@Router		/base/register/email [post]
func RegisterByEmail(ctx *gin.Context) {
	var req request.RegisterByEmail
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, service.RegisterByEmail(ctx, &req), nil)
}

// SendEmailCode
//
//	@Tags		guest
//	@Summary	邮箱注册验证码发送
//	@Accept		json
//	@Produce	application/json
//	@Param		data	body		request.SendEmailCode	true	"邮箱地址"
//	@Success	200		{object}	response.Response
//	@Router		/base/register/email/code [post]
func SendEmailCode(ctx *gin.Context) {
	var req request.SendEmailCode
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.SendEmailCode(ctx, req.Email)
	gin_util.Response(ctx, nil, err)
}
