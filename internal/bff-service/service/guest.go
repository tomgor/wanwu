package service

import (
	"fmt"
	"strconv"

	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/gin-gonic/gin"
)

const (
	customModeLight = "light"
	customModeDark  = "dark"
)

func GetLanguageSelect() *response.LanguageSelect {
	language := make([]response.Language, len(config.Cfg().I18n.Langs))
	for i, lang := range config.Cfg().I18n.Langs {
		language[i].Code = lang.Code
		language[i].Name = lang.Name
	}
	return &response.LanguageSelect{
		Languages:       language,
		DefaultLanguage: getLanguageByCode(config.Cfg().I18n.DefaultLang),
	}
}

func GetLogoCustomInfo(ctx *gin.Context, mode string) (response.LogoCustomInfo, error) {
	cfg := config.Cfg().CustomInfo
	ret := response.LogoCustomInfo{}
	var theme string
	switch mode {
	case customModeLight:
		theme = customModeLight
	case customModeDark:
		theme = customModeDark
	default:
		theme = cfg.DefaultMode
	}
	for _, mode := range cfg.Modes {
		if theme != mode.Mode {
			continue
		}
		ret = response.LogoCustomInfo{
			Login: response.CustomLogin{
				Background:       request.Avatar{Path: mode.Login.BackgroundPath},
				Logo:             request.Avatar{Path: mode.Login.LogoPath},
				LoginButtonColor: mode.Login.LoginButtonColor,
				WelcomeText:      gin_util.I18nKey(ctx, mode.Login.WelcomeText),
			},
			Home: response.CustomHome{
				Logo:            request.Avatar{Path: mode.Home.LogoPath},
				Title:           gin_util.I18nKey(ctx, mode.Home.Title),
				BackgroundColor: mode.Home.BackgroundColor,
			},
			Tab: response.CustomTab{
				Logo:  request.Avatar{Path: mode.Tab.TabLogoPath},
				Title: gin_util.I18nKey(ctx, mode.Tab.TabTitle),
			},
			About: response.CustomAbout{
				LogoPath:  mode.About.LogoPath,
				Version:   cfg.Version,
				Copyright: gin_util.I18nKey(ctx, mode.About.Copyright),
			},
			LinkList: config.Cfg().DocCenter.GetDocs(),
		}
		break
	}
	custom, err := operate.GetSystemCustom(ctx.Request.Context(), &operate_service.GetSystemCustomReq{Mode: theme})
	if err != nil {
		return ret, err
	}
	if custom.Tab.TabLogoPath != "" {
		ret.Tab.Logo = CacheAvatar(ctx, custom.Tab.TabLogoPath)
	}
	if custom.Tab.TabTitle != "" {
		ret.Tab.Title = custom.Tab.TabTitle
	}
	if custom.Login.LoginBgPath != "" {
		ret.Login.Background = CacheAvatar(ctx, custom.Login.LoginBgPath)
	}
	if custom.Login.LoginLogo != "" {
		ret.Login.Logo = CacheAvatar(ctx, custom.Login.LoginLogo)
	}
	if custom.Login.LoginButtonColor != "" {
		ret.Login.LoginButtonColor = custom.Login.LoginButtonColor
	}
	if custom.Login.LoginWelcomeText != "" {
		ret.Login.WelcomeText = custom.Login.LoginWelcomeText
	}
	if custom.Home.HomeName != "" {
		ret.Home.Title = custom.Home.HomeName
	}
	if custom.Home.HomeLogoPath != "" {
		ret.Home.Logo = CacheAvatar(ctx, custom.Home.HomeLogoPath)
	}
	if custom.Home.HomeBgColor != "" {
		ret.Home.BackgroundColor = custom.Home.HomeBgColor
	}
	return ret, nil
}

func GetCaptcha(ctx *gin.Context, key string) (*response.Captcha, error) {
	resp, err := iam.GetCaptcha(ctx.Request.Context(), &iam_service.GetCaptchaReq{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	return &response.Captcha{
		Key: key,
		B64: resp.B64,
	}, nil
}

func Login(ctx *gin.Context, login *request.Login, language string) (*response.Login, error) {
	password, err := decryptPD(login.Password)
	if err != nil {
		return nil, fmt.Errorf("decrypt password err: %v", err)
	}
	resp, err := iam.Login(ctx.Request.Context(), &iam_service.LoginReq{
		UserName: login.Username,
		Password: password,
		Key:      login.Key,
		Code:     login.Code,
		Language: language,
	})
	if err != nil {
		return nil, err
	}
	// orgs
	orgs, err := iam.GetOrgSelect(ctx.Request.Context(), &iam_service.GetOrgSelectReq{UserId: resp.User.GetUserId()})
	if err != nil {
		return nil, err
	}
	// jwt token
	claims, token, err := jwt_util.GenerateToken(
		resp.User.GetUserId(),
		jwt_util.UserTokenTimeout,
	)
	if err != nil {
		return nil, err
	}
	ctx.Set(gin_util.CLAIMS, &claims)
	// resp
	return &response.Login{
		UID:              resp.User.GetUserId(),
		Username:         resp.User.GetUserName(),
		Nickname:         resp.User.GetNickName(),
		Token:            token,
		ExpiresAt:        claims.StandardClaims.ExpiresAt * 1000, // 超时事件戳毫秒
		ExpireIn:         strconv.FormatInt(jwt_util.UserTokenTimeout, 10),
		Orgs:             toOrgIDNames(ctx, orgs.Selects, resp.User.GetUserId() == config.SystemAdminUserID),
		OrgPermission:    toOrgPermission(ctx, resp.Permission),
		Language:         getLanguageByCode(resp.User.Language),
		IsUpdatePassword: resp.Permission.LastUpdatePasswordAt != 0,
	}, nil
}

// --- internal ---

func getLanguageByCode(languageCode string) response.Language {
	langs := config.Cfg().I18n.Langs
	language := response.Language{Code: languageCode}
	for _, lang := range langs {
		if lang.Code == languageCode {
			language.Name = lang.Name
		}
	}
	return language
}

func toOrgPermission(ctx *gin.Context, orgPerm *iam_service.UserPermission) response.UserOrgPermission {
	return response.UserOrgPermission{
		IsAdmin:     orgPerm.IsAdmin,
		IsSystem:    orgPerm.IsSystem,
		Org:         toOrgIDName(ctx, orgPerm.Org),
		Roles:       toRoleIDNames(ctx, orgPerm.Roles),
		Permissions: toPermissions(orgPerm.IsAdmin, orgPerm.IsSystem, orgPerm.Perms),
	}
}

func toPermissions(isAdmin, isSystem bool, perms []*iam_service.Perm) []response.Permission {
	routes := mid.CollectPerms()
	var ret []response.Permission
	if isAdmin {
		for _, r := range routes {
			if isSystem && r.Tag == "permission.role" {
				continue
			}
			if !isSystem && r.Tag == "setting" {
				continue
			}
			ret = append(ret, response.Permission{
				Perm: r.Tag,
				Name: r.Name,
			})
		}
		return ret
	}
	for _, r := range routes {
		if r.Tag == "setting" {
			continue
		}
		for _, perm := range perms {
			if perm.Perm == r.Tag {
				ret = append(ret, response.Permission{
					Perm: r.Tag,
					Name: r.Name,
				})
				break
			}
		}
	}
	return ret
}
